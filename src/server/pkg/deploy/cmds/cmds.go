package cmds

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pachyderm/pachyderm/src/client"
	"github.com/pachyderm/pachyderm/src/client/auth"
	"github.com/pachyderm/pachyderm/src/client/enterprise"
	"github.com/pachyderm/pachyderm/src/client/pkg/config"
	"github.com/pachyderm/pachyderm/src/client/pkg/errors"
	"github.com/pachyderm/pachyderm/src/client/pkg/grpcutil"
	"github.com/pachyderm/pachyderm/src/client/pkg/helm"
	"github.com/pachyderm/pachyderm/src/client/version"
	"github.com/pachyderm/pachyderm/src/server/pkg/cmdutil"
	"github.com/pachyderm/pachyderm/src/server/pkg/deploy"
	"github.com/pachyderm/pachyderm/src/server/pkg/deploy/assets"
	"github.com/pachyderm/pachyderm/src/server/pkg/deploy/images"
	_metrics "github.com/pachyderm/pachyderm/src/server/pkg/metrics"
	"github.com/pachyderm/pachyderm/src/server/pkg/obj"
	"github.com/pachyderm/pachyderm/src/server/pkg/serde"
	clientcmd "k8s.io/client-go/tools/clientcmd/api/v1"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	awsAccessKeyIDRE = regexp.MustCompile("^[A-Z0-9]{20}$")
	awsSecretRE      = regexp.MustCompile("^[A-Za-z0-9/+=]{40}$")
	awsRegionRE      = regexp.MustCompile("^[a-z]{2}(?:-gov)?-[a-z]+-[0-9]$")
)

const (
	defaultDashImage       = "pachyderm/dash:0.5.48"
	defaultIDEHubImage     = "pachyderm/ide-hub:1.0.0"
	defaultIDEUserImage    = "pachyderm/ide-user:1.0.0"
	jupyterhubChartVersion = "0.8.2"
)

func kubectl(stdin io.Reader, context *config.Context, args ...string) error {
	var environ []string = nil
	if context != nil {
		tmpfile, err := ioutil.TempFile("", "transient-kube-config-*.yaml")
		if err != nil {
			return errors.Wrapf(err, "failed to create transient kube config")
		}
		defer os.Remove(tmpfile.Name())

		config := clientcmd.Config{
			Kind:           "Config",
			APIVersion:     "v1",
			CurrentContext: "pachyderm-active-context",
			Contexts: []clientcmd.NamedContext{
				clientcmd.NamedContext{
					Name: "pachyderm-active-context",
					Context: clientcmd.Context{
						Cluster:   context.ClusterName,
						AuthInfo:  context.AuthInfo,
						Namespace: context.Namespace,
					},
				},
			},
		}

		var buf bytes.Buffer
		if err := encoder("yaml", &buf).Encode(config); err != nil {
			return errors.Wrapf(err, "failed to encode config")
		}

		tmpfile.Write(buf.Bytes())
		tmpfile.Close()

		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return errors.Wrapf(err, "failed to discover default kube config: could not get user home directory")
			}
			kubeconfig = path.Join(home, ".kube", "config")
			if _, err = os.Stat(kubeconfig); os.IsNotExist(err) {
				return errors.Wrapf(err, "failed to discover default kube config: %q does not exist", kubeconfig)
			}
		}
		kubeconfig = fmt.Sprintf("%s%c%s", kubeconfig, os.PathListSeparator, tmpfile.Name())

		// note that this will override `KUBECONFIG` (if it is already defined) in
		// the environment; see examples under
		// https://golang.org/pkg/os/exec/#Command
		environ = os.Environ()
		environ = append(environ, fmt.Sprintf("KUBECONFIG=%s", kubeconfig))

		if stdin == nil {
			stdin = os.Stdin
		}
	}

	ioObj := cmdutil.IO{
		Stdin:   stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Environ: environ,
	}

	args = append([]string{"kubectl"}, args...)
	return cmdutil.RunIO(ioObj, args...)
}

// Generates a random secure token, in hex
func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// Return the appropriate encoder for the given output format.
func encoder(output string, w io.Writer) serde.Encoder {
	if output == "" {
		output = "json"
	} else {
		output = strings.ToLower(output)
	}
	e, err := serde.GetEncoder(output, w,
		serde.WithIndent(2),
		serde.WithOrigName(true),
	)
	if err != nil {
		cmdutil.ErrorAndExit(err.Error())
	}
	return e
}

func kubectlCreate(dryRun bool, manifest []byte, opts *assets.AssetOpts) error {
	if dryRun {
		_, err := os.Stdout.Write(manifest)
		return err
	}
	// we set --validate=false due to https://github.com/kubernetes/kubernetes/issues/53309
	if err := kubectl(bytes.NewReader(manifest), nil, "apply", "-f", "-", "--validate=false", "--namespace", opts.Namespace); err != nil {
		return err
	}

	fmt.Println("\nPachyderm is launching. Check its status with \"kubectl get all\"")
	if opts.DashOnly || !opts.NoDash {
		fmt.Println("Once launched, access the dashboard by running \"pachctl port-forward\"")
	}
	fmt.Println("")

	return nil
}

// findEquivalentContext searches for a context in the existing config that
// references the same cluster as the context passed in. If no such context
// was found, default values are returned instead.
func findEquivalentContext(cfg *config.Config, to *config.Context) (string, *config.Context) {
	// first check the active context
	activeContextName, activeContext, _ := cfg.ActiveContext(false)
	if activeContextName != "" && to.EqualClusterReference(activeContext) {
		return activeContextName, activeContext
	}

	// failing that, search all contexts (sorted by name to be deterministic)
	contextNames := []string{}
	for contextName := range cfg.V2.Contexts {
		contextNames = append(contextNames, contextName)
	}
	sort.Strings(contextNames)
	for _, contextName := range contextNames {
		existingContext := cfg.V2.Contexts[contextName]

		if to.EqualClusterReference(existingContext) {
			return contextName, existingContext
		}
	}

	return "", nil
}

func contextCreate(namePrefix, namespace, serverCert string) error {
	kubeConfig, err := config.RawKubeConfig()
	if err != nil {
		return err
	}
	kubeContext := kubeConfig.Contexts[kubeConfig.CurrentContext]

	clusterName := ""
	authInfo := ""
	if kubeContext != nil {
		clusterName = kubeContext.Cluster
		authInfo = kubeContext.AuthInfo
	}

	cfg, err := config.Read(false)
	if err != nil {
		return err
	}

	newContext := &config.Context{
		Source:      config.ContextSource_IMPORTED,
		ClusterName: clusterName,
		AuthInfo:    authInfo,
		Namespace:   namespace,
		ServerCAs:   serverCert,
	}

	equivalentContextName, equivalentContext := findEquivalentContext(cfg, newContext)
	if equivalentContext != nil {
		cfg.V2.ActiveContext = equivalentContextName
		equivalentContext.Source = newContext.Source
		equivalentContext.ClusterDeploymentID = ""
		equivalentContext.ServerCAs = newContext.ServerCAs
		return cfg.Write()
	}

	// we couldn't find an existing context that is the same as the new one,
	// so we'll have to create it
	newContextName := namePrefix
	if _, ok := cfg.V2.Contexts[newContextName]; ok {
		newContextName = fmt.Sprintf("%s-%s", namePrefix, time.Now().Format("2006-01-02-15-04-05"))
	}

	cfg.V2.Contexts[newContextName] = newContext
	cfg.V2.ActiveContext = newContextName
	return cfg.Write()
}

// containsEmpty is a helper function used for validation (particularly for
// validating that creds arguments aren't empty
func containsEmpty(vals []string) bool {
	for _, val := range vals {
		if val == "" {
			return true
		}
	}
	return false
}

func standardDeployCmds() []*cobra.Command {
	var commands []*cobra.Command
	var opts *assets.AssetOpts

	var dryRun bool
	var outputFormat string
	var namespace string
	var serverCert string
	var blockCacheSize string
	var dashImage string
	var dashOnly bool
	var etcdCPURequest string
	var etcdMemRequest string
	var etcdNodes int
	var etcdStorageClassName string
	var etcdVolume string
	var exposeObjectAPI bool
	var imagePullSecret string
	var localRoles bool
	var logLevel string
	var newStorageLayer bool
	var noDash bool
	var noExposeDockerSocket bool
	var noGuaranteed bool
	var noRBAC bool
	var pachdCPURequest string
	var pachdNonCacheMemRequest string
	var pachdShards int
	var registry string
	var tlsCertKey string
	var uploadConcurrencyLimit int
	var putFileConcurrencyLimit int
	var clusterDeploymentID string
	var requireCriticalServersOnly bool
	var workerServiceAccountName string
	appendGlobalFlags := func(cmd *cobra.Command) {
		cmd.Flags().IntVar(&pachdShards, "shards", 16, "(rarely set) The maximum number of pachd nodes allowed in the cluster; increasing this number blindly can result in degraded performance.")
		cmd.Flags().IntVar(&etcdNodes, "dynamic-etcd-nodes", 0, "Deploy etcd as a StatefulSet with the given number of pods.  The persistent volumes used by these pods are provisioned dynamically.  Note that StatefulSet is currently a beta kubernetes feature, which might be unavailable in older versions of kubernetes.")
		cmd.Flags().StringVar(&etcdVolume, "static-etcd-volume", "", "Deploy etcd as a ReplicationController with one pod.  The pod uses the given persistent volume.")
		cmd.Flags().StringVar(&etcdStorageClassName, "etcd-storage-class", "", "If set, the name of an existing StorageClass to use for etcd storage. Ignored if --static-etcd-volume is set.")
		cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Don't actually deploy pachyderm to Kubernetes, instead just print the manifest. Note that a pachyderm context will not be created, unless you also use `--create-context`.")
		cmd.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format. One of: json|yaml")
		cmd.Flags().StringVar(&logLevel, "log-level", "info", "The level of log messages to print options are, from least to most verbose: \"error\", \"info\", \"debug\".")
		cmd.Flags().BoolVar(&dashOnly, "dashboard-only", false, "Only deploy the Pachyderm UI (experimental), without the rest of pachyderm. This is for launching the UI adjacent to an existing Pachyderm cluster. After deployment, run \"pachctl port-forward\" to connect")
		cmd.Flags().BoolVar(&noDash, "no-dashboard", false, "Don't deploy the Pachyderm UI alongside Pachyderm (experimental).")
		cmd.Flags().StringVar(&registry, "registry", "", "The registry to pull images from.")
		cmd.Flags().StringVar(&imagePullSecret, "image-pull-secret", "", "A secret in Kubernetes that's needed to pull from your private registry.")
		cmd.Flags().StringVar(&dashImage, "dash-image", "", "Image URL for pachyderm dashboard")
		cmd.Flags().BoolVar(&noGuaranteed, "no-guaranteed", false, "Don't use guaranteed QoS for etcd and pachd deployments. Turning this on (turning guaranteed QoS off) can lead to more stable local clusters (such as on Minikube), it should normally be used for production clusters.")
		cmd.Flags().BoolVar(&noRBAC, "no-rbac", false, "Don't deploy RBAC roles for Pachyderm. (for k8s versions prior to 1.8)")
		cmd.Flags().BoolVar(&localRoles, "local-roles", false, "Use namespace-local roles instead of cluster roles. Ignored if --no-rbac is set.")
		cmd.Flags().StringVar(&namespace, "namespace", "", "Kubernetes namespace to deploy Pachyderm to.")
		cmd.Flags().BoolVar(&noExposeDockerSocket, "no-expose-docker-socket", false, "Don't expose the Docker socket to worker containers. This limits the privileges of workers which prevents them from automatically setting the container's working dir and user.")
		cmd.Flags().BoolVar(&exposeObjectAPI, "expose-object-api", false, "If set, instruct pachd to serve its object/block API on its public port (not safe with auth enabled, do not set in production).")
		cmd.Flags().StringVar(&tlsCertKey, "tls", "", "string of the form \"<cert path>,<key path>\" of the signed TLS certificate and private key that Pachd should use for TLS authentication (enables TLS-encrypted communication with Pachd)")
		cmd.Flags().BoolVar(&newStorageLayer, "new-storage-layer", false, "(feature flag) Do not set, used for testing.")
		cmd.Flags().IntVar(&uploadConcurrencyLimit, "upload-concurrency-limit", assets.DefaultUploadConcurrencyLimit, "The maximum number of concurrent object storage uploads per Pachd instance.")
		cmd.Flags().IntVar(&putFileConcurrencyLimit, "put-file-concurrency-limit", assets.DefaultPutFileConcurrencyLimit, "The maximum number of files to upload or fetch from remote sources (HTTP, blob storage) using PutFile concurrently.")
		cmd.Flags().StringVar(&clusterDeploymentID, "cluster-deployment-id", "", "Set an ID for the cluster deployment. Defaults to a random value.")
		cmd.Flags().BoolVar(&requireCriticalServersOnly, "require-critical-servers-only", assets.DefaultRequireCriticalServersOnly, "Only require the critical Pachd servers to startup and run without errors.")
		cmd.Flags().StringVar(&workerServiceAccountName, "worker-service-account", assets.DefaultWorkerServiceAccountName, "The Kubernetes service account for workers to use when creating S3 gateways.")

		// Flags for setting pachd resource requests. These should rarely be set --
		// only if we get the defaults wrong, or users have an unusual access pattern
		//
		// All of these are empty by default, because the actual default values depend
		// on the backend to which we're. The defaults are set in
		// s/s/pkg/deploy/assets/assets.go
		cmd.Flags().StringVar(&pachdCPURequest,
			"pachd-cpu-request", "", "(rarely set) The size of Pachd's CPU "+
				"request, which we give to Kubernetes. Size is in cores (with partial "+
				"cores allowed and encouraged).")
		cmd.Flags().StringVar(&blockCacheSize, "block-cache-size", "",
			"Size of pachd's in-memory cache for PFS files. Size is specified in "+
				"bytes, with allowed SI suffixes (M, K, G, Mi, Ki, Gi, etc).")
		cmd.Flags().StringVar(&pachdNonCacheMemRequest,
			"pachd-memory-request", "", "(rarely set) The size of PachD's memory "+
				"request in addition to its block cache (set via --block-cache-size). "+
				"Size is in bytes, with SI suffixes (M, K, G, Mi, Ki, Gi, etc).")
		cmd.Flags().StringVar(&etcdCPURequest,
			"etcd-cpu-request", "", "(rarely set) The size of etcd's CPU request, "+
				"which we give to Kubernetes. Size is in cores (with partial cores "+
				"allowed and encouraged).")
		cmd.Flags().StringVar(&etcdMemRequest,
			"etcd-memory-request", "", "(rarely set) The size of etcd's memory "+
				"request. Size is in bytes, with SI suffixes (M, K, G, Mi, Ki, Gi, "+
				"etc).")
	}

	var retries int
	var timeout string
	var uploadACL string
	var reverse bool
	var partSize int64
	var maxUploadParts int
	var disableSSL bool
	var noVerifySSL bool
	appendS3Flags := func(cmd *cobra.Command) {
		cmd.Flags().IntVar(&retries, "retries", obj.DefaultRetries, "(rarely set) Set a custom number of retries for object storage requests.")
		cmd.Flags().StringVar(&timeout, "timeout", obj.DefaultTimeout, "(rarely set) Set a custom timeout for object storage requests.")
		cmd.Flags().StringVar(&uploadACL, "upload-acl", obj.DefaultUploadACL, "(rarely set) Set a custom upload ACL for object storage uploads.")
		cmd.Flags().BoolVar(&reverse, "reverse", obj.DefaultReverse, "(rarely set) Reverse object storage paths.")
		cmd.Flags().Int64Var(&partSize, "part-size", obj.DefaultPartSize, "(rarely set) Set a custom part size for object storage uploads.")
		cmd.Flags().IntVar(&maxUploadParts, "max-upload-parts", obj.DefaultMaxUploadParts, "(rarely set) Set a custom maximum number of upload parts.")
		cmd.Flags().BoolVar(&disableSSL, "disable-ssl", obj.DefaultDisableSSL, "(rarely set) Disable SSL.")
		cmd.Flags().BoolVar(&noVerifySSL, "no-verify-ssl", obj.DefaultNoVerifySSL, "(rarely set) Skip SSL certificate verification (typically used for enabling self-signed certificates).")
	}

	var contextName string
	var createContext bool
	appendContextFlags := func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&contextName, "context", "c", "", "Name of the context to add to the pachyderm config. If unspecified, a context name will automatically be derived.")
		cmd.Flags().BoolVar(&createContext, "create-context", false, "Create a context, even with `--dry-run`.")
	}

	preRunInternal := func(args []string) error {
		cfg, err := config.Read(false)
		if err != nil {
			log.Warningf("could not read config to check whether cluster metrics "+
				"will be enabled: %v.\n", err)
		}

		if namespace == "" {
			kubeConfig := config.KubeConfig(nil)
			var err error
			namespace, _, err = kubeConfig.Namespace()
			if err != nil {
				log.Warningf("using namespace \"default\" (couldn't load namespace "+
					"from kubernetes config: %v)\n", err)
				namespace = "default"
			}
		}

		dashImage = getDefaultOrLatestDashImage(dashImage, dryRun)
		opts = &assets.AssetOpts{
			FeatureFlags: assets.FeatureFlags{
				NewStorageLayer: newStorageLayer,
			},
			StorageOpts: assets.StorageOpts{
				UploadConcurrencyLimit:  uploadConcurrencyLimit,
				PutFileConcurrencyLimit: putFileConcurrencyLimit,
			},
			PachdShards:                uint64(pachdShards),
			Version:                    version.PrettyPrintVersion(version.Version),
			LogLevel:                   logLevel,
			Metrics:                    cfg == nil || cfg.V2.Metrics,
			PachdCPURequest:            pachdCPURequest,
			PachdNonCacheMemRequest:    pachdNonCacheMemRequest,
			BlockCacheSize:             blockCacheSize,
			EtcdCPURequest:             etcdCPURequest,
			EtcdMemRequest:             etcdMemRequest,
			EtcdNodes:                  etcdNodes,
			EtcdVolume:                 etcdVolume,
			EtcdStorageClassName:       etcdStorageClassName,
			DashOnly:                   dashOnly,
			NoDash:                     noDash,
			DashImage:                  dashImage,
			Registry:                   registry,
			ImagePullSecret:            imagePullSecret,
			NoGuaranteed:               noGuaranteed,
			NoRBAC:                     noRBAC,
			LocalRoles:                 localRoles,
			Namespace:                  namespace,
			NoExposeDockerSocket:       noExposeDockerSocket,
			ExposeObjectAPI:            exposeObjectAPI,
			ClusterDeploymentID:        clusterDeploymentID,
			RequireCriticalServersOnly: requireCriticalServersOnly,
			WorkerServiceAccountName:   workerServiceAccountName,
		}
		if tlsCertKey != "" {
			// TODO(msteffen): If either the cert path or the key path contains a
			// comma, this doesn't work
			certKey := strings.Split(tlsCertKey, ",")
			if len(certKey) != 2 {
				return fmt.Errorf("could not split TLS certificate and key correctly; must have two parts but got: %#v", certKey)
			}
			opts.TLS = &assets.TLSOpts{
				ServerCert: certKey[0],
				ServerKey:  certKey[1],
			}

			serverCertBytes, err := ioutil.ReadFile(certKey[0])
			if err != nil {
				return errors.Wrapf(err, "could not read server cert at %q", certKey[0])
			}
			serverCert = base64.StdEncoding.EncodeToString([]byte(serverCertBytes))
		}
		return nil
	}
	preRun := cmdutil.Run(preRunInternal)

	deployPreRun := cmdutil.Run(func(args []string) error {
		if version.IsUnstable() {
			fmt.Printf("WARNING: The version of Pachyderm you are deploying (%s) is an unstable pre-release build and may not support data migration.\n\n", version.PrettyVersion())

			if ok, err := cmdutil.InteractiveConfirm(); err != nil {
				return err
			} else if !ok {
				return errors.New("deploy aborted")
			}
		}
		return preRunInternal(args)
	})

	var dev bool
	var hostPath string
	deployLocal := &cobra.Command{
		Short:  "Deploy a single-node Pachyderm cluster with local metadata storage.",
		Long:   "Deploy a single-node Pachyderm cluster with local metadata storage.",
		PreRun: deployPreRun,
		Run: cmdutil.RunFixedArgs(0, func(args []string) (retErr error) {
			if !dev {
				start := time.Now()
				startMetricsWait := _metrics.StartReportAndFlushUserAction("Deploy", start)
				defer startMetricsWait()
				defer func() {
					finishMetricsWait := _metrics.FinishReportAndFlushUserAction("Deploy", retErr, start)
					finishMetricsWait()
				}()
			}
			if dev {
				// Use dev build instead of release build
				opts.Version = deploy.DevVersionTag

				// we turn metrics off if this is a dev cluster. The default
				// is set by deploy.PersistentPreRun, below.
				opts.Metrics = false

				// Disable authentication, for tests
				opts.DisableAuthentication = true

				// Serve the Pachyderm object/block API locally, as this is needed by
				// our tests (and authentication is disabled anyway)
				opts.ExposeObjectAPI = true
			}
			var buf bytes.Buffer
			if err := assets.WriteLocalAssets(
				encoder(outputFormat, &buf), opts, hostPath,
			); err != nil {
				return err
			}
			if err := kubectlCreate(dryRun, buf.Bytes(), opts); err != nil {
				return err
			}
			if !dryRun || createContext {
				if contextName == "" {
					contextName = "local"
				}
				if err := contextCreate(contextName, namespace, serverCert); err != nil {
					return err
				}
			}
			return nil
		}),
	}
	appendGlobalFlags(deployLocal)
	appendContextFlags(deployLocal)
	deployLocal.Flags().StringVar(&hostPath, "host-path", "/var/pachyderm", "Location on the host machine where PFS metadata will be stored.")
	deployLocal.Flags().BoolVarP(&dev, "dev", "d", false, "Deploy pachd with local version tags, disable metrics, expose Pachyderm's object/block API, and use an insecure authentication mechanism (do not set on any cluster with sensitive data)")
	commands = append(commands, cmdutil.CreateAlias(deployLocal, "deploy local"))

	deployGoogle := &cobra.Command{
		Use:   "{{alias}} <bucket-name> <disk-size> [<credentials-file>]",
		Short: "Deploy a Pachyderm cluster running on Google Cloud Platform.",
		Long: `Deploy a Pachyderm cluster running on Google Cloud Platform.
  <bucket-name>: A Google Cloud Storage bucket where Pachyderm will store PFS data.
  <disk-size>: Size of Google Compute Engine persistent disks in GB (assumed to all be the same).
  <credentials-file>: A file containing the private key for the account (downloaded from Google Compute Engine).`,
		PreRun: deployPreRun,
		Run: cmdutil.RunBoundedArgs(2, 3, func(args []string) (retErr error) {
			start := time.Now()
			startMetricsWait := _metrics.StartReportAndFlushUserAction("Deploy", start)
			defer startMetricsWait()
			defer func() {
				finishMetricsWait := _metrics.FinishReportAndFlushUserAction("Deploy", retErr, start)
				finishMetricsWait()
			}()
			volumeSize, err := strconv.Atoi(args[1])
			if err != nil {
				return errors.Errorf("volume size needs to be an integer; instead got %v", args[1])
			}
			var buf bytes.Buffer
			opts.BlockCacheSize = "0G" // GCS is fast so we want to disable the block cache. See issue #1650
			var cred string
			if len(args) == 3 {
				credBytes, err := ioutil.ReadFile(args[2])
				if err != nil {
					return errors.Wrapf(err, "error reading creds file %s", args[2])
				}
				cred = string(credBytes)
			}
			bucket := strings.TrimPrefix(args[0], "gs://")
			if err = assets.WriteGoogleAssets(
				encoder(outputFormat, &buf), opts, bucket, cred, volumeSize,
			); err != nil {
				return err
			}
			if err := kubectlCreate(dryRun, buf.Bytes(), opts); err != nil {
				return err
			}
			if !dryRun || createContext {
				if contextName == "" {
					contextName = "gcs"
				}
				if err := contextCreate(contextName, namespace, serverCert); err != nil {
					return err
				}
			}
			return nil
		}),
	}
	appendGlobalFlags(deployGoogle)
	appendContextFlags(deployGoogle)
	commands = append(commands, cmdutil.CreateAlias(deployGoogle, "deploy google"))

	var objectStoreBackend string
	var persistentDiskBackend string
	var secure bool
	var isS3V2 bool
	deployCustom := &cobra.Command{
		Use:   "{{alias}} --persistent-disk <persistent disk backend> --object-store <object store backend> <persistent disk args> <object store args>",
		Short: "Deploy a custom Pachyderm cluster configuration",
		Long: `Deploy a custom Pachyderm cluster configuration.
If <object store backend> is \"s3\", then the arguments are:
    <volumes> <size of volumes (in GB)> <bucket> <id> <secret> <endpoint>`,
		PreRun: deployPreRun,
		Run: cmdutil.RunBoundedArgs(4, 7, func(args []string) (retErr error) {
			start := time.Now()
			startMetricsWait := _metrics.StartReportAndFlushUserAction("Deploy", start)
			defer startMetricsWait()
			defer func() {
				finishMetricsWait := _metrics.FinishReportAndFlushUserAction("Deploy", retErr, start)
				finishMetricsWait()
			}()
			// Setup advanced configuration.
			advancedConfig := &obj.AmazonAdvancedConfiguration{
				Retries:        retries,
				Timeout:        timeout,
				UploadACL:      uploadACL,
				Reverse:        reverse,
				PartSize:       partSize,
				MaxUploadParts: maxUploadParts,
				DisableSSL:     disableSSL,
				NoVerifySSL:    noVerifySSL,
			}
			// Generate manifest and write assets.
			var buf bytes.Buffer
			if err := assets.WriteCustomAssets(
				encoder(outputFormat, &buf), opts, args, objectStoreBackend,
				persistentDiskBackend, secure, isS3V2, advancedConfig,
			); err != nil {
				return err
			}
			if err := kubectlCreate(dryRun, buf.Bytes(), opts); err != nil {
				return err
			}
			if !dryRun || createContext {
				if contextName == "" {
					contextName = "custom"
				}
				if err := contextCreate(contextName, namespace, serverCert); err != nil {
					return err
				}
			}
			return nil
		}),
	}
	appendGlobalFlags(deployCustom)
	appendS3Flags(deployCustom)
	appendContextFlags(deployCustom)
	// (bryce) secure should be merged with disableSSL, but it would be a breaking change.
	deployCustom.Flags().BoolVarP(&secure, "secure", "s", false, "Enable secure access to a Minio server.")
	deployCustom.Flags().StringVar(&persistentDiskBackend, "persistent-disk", "aws",
		"(required) Backend providing persistent local volumes to stateful pods. "+
			"One of: aws, google, or azure.")
	deployCustom.Flags().StringVar(&objectStoreBackend, "object-store", "s3",
		"(required) Backend providing an object-storage API to pachyderm. One of: "+
			"s3, gcs, or azure-blob.")
	deployCustom.Flags().BoolVar(&isS3V2, "isS3V2", false, "Enable S3V2 client")
	commands = append(commands, cmdutil.CreateAlias(deployCustom, "deploy custom"))

	var cloudfrontDistribution string
	var creds string
	var iamRole string
	var vault string
	deployAmazon := &cobra.Command{
		Use:   "{{alias}} <bucket-name> <region> <disk-size>",
		Short: "Deploy a Pachyderm cluster running on AWS.",
		Long: `Deploy a Pachyderm cluster running on AWS.
  <bucket-name>: An S3 bucket where Pachyderm will store PFS data.
  <region>: The AWS region where Pachyderm is being deployed (e.g. us-west-1)
  <disk-size>: Size of EBS volumes, in GB (assumed to all be the same).`,
		PreRun: deployPreRun,
		Run: cmdutil.RunFixedArgs(3, func(args []string) (retErr error) {
			start := time.Now()
			startMetricsWait := _metrics.StartReportAndFlushUserAction("Deploy", start)
			defer startMetricsWait()
			defer func() {
				finishMetricsWait := _metrics.FinishReportAndFlushUserAction("Deploy", retErr, start)
				finishMetricsWait()
			}()
			if creds == "" && vault == "" && iamRole == "" {
				return errors.Errorf("one of --credentials, --vault, or --iam-role needs to be provided")
			}

			// populate 'amazonCreds' & validate
			var amazonCreds *assets.AmazonCreds
			s := bufio.NewScanner(os.Stdin)
			if creds != "" {
				parts := strings.Split(creds, ",")
				if len(parts) < 2 || len(parts) > 3 || containsEmpty(parts[:2]) {
					return errors.Errorf("incorrect format of --credentials")
				}
				amazonCreds = &assets.AmazonCreds{ID: parts[0], Secret: parts[1]}
				if len(parts) > 2 {
					amazonCreds.Token = parts[2]
				}

				if !awsAccessKeyIDRE.MatchString(amazonCreds.ID) {
					fmt.Fprintf(os.Stderr, "The AWS Access Key seems invalid (does not "+
						"match %q). Do you want to continue deploying? [yN]\n",
						awsAccessKeyIDRE)
					if s.Scan(); s.Text()[0] != 'y' && s.Text()[0] != 'Y' {
						os.Exit(1)
					}
				}

				if !awsSecretRE.MatchString(amazonCreds.Secret) {
					fmt.Fprintf(os.Stderr, "The AWS Secret seems invalid (does not "+
						"match %q). Do you want to continue deploying? [yN]\n", awsSecretRE)
					if s.Scan(); s.Text()[0] != 'y' && s.Text()[0] != 'Y' {
						os.Exit(1)
					}
				}
			}
			if vault != "" {
				if amazonCreds != nil {
					return errors.Errorf("only one of --credentials, --vault, or --iam-role needs to be provided")
				}
				parts := strings.Split(vault, ",")
				if len(parts) != 3 || containsEmpty(parts) {
					return errors.Errorf("incorrect format of --vault")
				}
				amazonCreds = &assets.AmazonCreds{VaultAddress: parts[0], VaultRole: parts[1], VaultToken: parts[2]}
			}
			if iamRole != "" {
				if amazonCreds != nil {
					return errors.Errorf("only one of --credentials, --vault, or --iam-role needs to be provided")
				}
				opts.IAMRole = iamRole
			}
			volumeSize, err := strconv.Atoi(args[2])
			if err != nil {
				return errors.Errorf("volume size needs to be an integer; instead got %v", args[2])
			}
			if strings.TrimSpace(cloudfrontDistribution) != "" {
				log.Warningf("you specified a cloudfront distribution; deploying on " +
					"AWS with cloudfront is currently an alpha feature. No security " +
					"restrictions have been applied to cloudfront, making all data " +
					"public (obscured but not secured)\n")
			}
			bucket, region := strings.TrimPrefix(args[0], "s3://"), args[1]
			if !awsRegionRE.MatchString(region) {
				fmt.Fprintf(os.Stderr, "The AWS region seems invalid (does not match "+
					"%q). Do you want to continue deploying? [yN]\n", awsRegionRE)
				if s.Scan(); s.Text()[0] != 'y' && s.Text()[0] != 'Y' {
					os.Exit(1)
				}
			}
			// Setup advanced configuration.
			advancedConfig := &obj.AmazonAdvancedConfiguration{
				Retries:        retries,
				Timeout:        timeout,
				UploadACL:      uploadACL,
				Reverse:        reverse,
				PartSize:       partSize,
				MaxUploadParts: maxUploadParts,
				DisableSSL:     disableSSL,
				NoVerifySSL:    noVerifySSL,
			}
			// Generate manifest and write assets.
			var buf bytes.Buffer
			if err = assets.WriteAmazonAssets(
				encoder(outputFormat, &buf), opts, region, bucket, volumeSize,
				amazonCreds, cloudfrontDistribution, advancedConfig,
			); err != nil {
				return err
			}
			if err := kubectlCreate(dryRun, buf.Bytes(), opts); err != nil {
				return err
			}
			if !dryRun || createContext {
				if contextName == "" {
					contextName = "aws"
				}
				if err := contextCreate(contextName, namespace, serverCert); err != nil {
					return err
				}
			}
			return nil
		}),
	}
	appendGlobalFlags(deployAmazon)
	appendS3Flags(deployAmazon)
	appendContextFlags(deployAmazon)
	deployAmazon.Flags().StringVar(&cloudfrontDistribution, "cloudfront-distribution", "",
		"Deploying on AWS with cloudfront is currently "+
			"an alpha feature. No security restrictions have been"+
			"applied to cloudfront, making all data public (obscured but not secured)")
	deployAmazon.Flags().StringVar(&creds, "credentials", "", "Use the format \"<id>,<secret>[,<token>]\". You can get a token by running \"aws sts get-session-token\".")
	deployAmazon.Flags().StringVar(&vault, "vault", "", "Use the format \"<address/hostport>,<role>,<token>\".")
	deployAmazon.Flags().StringVar(&iamRole, "iam-role", "", fmt.Sprintf("Use the given IAM role for authorization, as opposed to using static credentials. The given role will be applied as the annotation %s, this used with a Kubernetes IAM role management system such as kube2iam allows you to give pachd credentials in a more secure way.", assets.IAMAnnotation))
	commands = append(commands, cmdutil.CreateAlias(deployAmazon, "deploy amazon"))

	deployMicrosoft := &cobra.Command{
		Use:   "{{alias}} <container> <account-name> <account-key> <disk-size>",
		Short: "Deploy a Pachyderm cluster running on Microsoft Azure.",
		Long: `Deploy a Pachyderm cluster running on Microsoft Azure.
  <container>: An Azure container where Pachyderm will store PFS data.
  <disk-size>: Size of persistent volumes, in GB (assumed to all be the same).`,
		PreRun: deployPreRun,
		Run: cmdutil.RunFixedArgs(4, func(args []string) (retErr error) {
			start := time.Now()
			startMetricsWait := _metrics.StartReportAndFlushUserAction("Deploy", start)
			defer startMetricsWait()
			defer func() {
				finishMetricsWait := _metrics.FinishReportAndFlushUserAction("Deploy", retErr, start)
				finishMetricsWait()
			}()
			if _, err := base64.StdEncoding.DecodeString(args[2]); err != nil {
				return errors.Errorf("storage-account-key needs to be base64 encoded; instead got '%v'", args[2])
			}
			if opts.EtcdVolume != "" {
				tempURI, err := url.ParseRequestURI(opts.EtcdVolume)
				if err != nil {
					return errors.Errorf("volume URI needs to be a well-formed URI; instead got '%v'", opts.EtcdVolume)
				}
				opts.EtcdVolume = tempURI.String()
			}
			volumeSize, err := strconv.Atoi(args[3])
			if err != nil {
				return errors.Errorf("volume size needs to be an integer; instead got %v", args[3])
			}
			var buf bytes.Buffer
			container := strings.TrimPrefix(args[0], "wasb://")
			accountName, accountKey := args[1], args[2]
			if err = assets.WriteMicrosoftAssets(
				encoder(outputFormat, &buf), opts, container, accountName, accountKey, volumeSize,
			); err != nil {
				return err
			}
			if err := kubectlCreate(dryRun, buf.Bytes(), opts); err != nil {
				return err
			}
			if !dryRun || createContext {
				if contextName == "" {
					contextName = "azure"
				}
				if err := contextCreate(contextName, namespace, serverCert); err != nil {
					return err
				}
			}
			return nil
		}),
	}
	appendGlobalFlags(deployMicrosoft)
	appendContextFlags(deployMicrosoft)
	commands = append(commands, cmdutil.CreateAlias(deployMicrosoft, "deploy microsoft"))

	deployStorageSecrets := func(data map[string][]byte) error {
		cfg, err := config.Read(false)
		if err != nil {
			return err
		}
		_, activeContext, err := cfg.ActiveContext(true)
		if err != nil {
			return err
		}

		// clean up any empty, but non-nil strings in the data, since those will prevent those fields from getting merged when we do the patch
		for k, v := range data {
			if v != nil && len(v) == 0 {
				delete(data, k)
			}
		}

		var buf bytes.Buffer
		if err = assets.WriteSecret(encoder(outputFormat, &buf), data, opts); err != nil {
			return err
		}
		if dryRun {
			_, err := os.Stdout.Write(buf.Bytes())
			return err
		}

		s := buf.String()
		return kubectl(&buf, activeContext, "patch", "secret", "pachyderm-storage-secret", "-p", s, "--namespace", opts.Namespace, "--type=merge")
	}

	deployStorageAmazon := &cobra.Command{
		Use:    "{{alias}} <region> <access-key-id> <secret-access-key> [<session-token>]",
		Short:  "Deploy credentials for the Amazon S3 storage provider.",
		Long:   "Deploy credentials for the Amazon S3 storage provider, so that Pachyderm can ingress data from and egress data to it.",
		PreRun: preRun,
		Run: cmdutil.RunBoundedArgs(3, 4, func(args []string) error {
			var token string
			if len(args) == 4 {
				token = args[3]
			}
			// Setup advanced configuration.
			advancedConfig := &obj.AmazonAdvancedConfiguration{
				Retries:        retries,
				Timeout:        timeout,
				UploadACL:      uploadACL,
				Reverse:        reverse,
				PartSize:       partSize,
				MaxUploadParts: maxUploadParts,
				DisableSSL:     disableSSL,
				NoVerifySSL:    noVerifySSL,
			}
			return deployStorageSecrets(assets.AmazonSecret(args[0], "", args[1], args[2], token, "", "", advancedConfig))
		}),
	}
	appendGlobalFlags(deployStorageAmazon)
	appendS3Flags(deployStorageAmazon)
	commands = append(commands, cmdutil.CreateAlias(deployStorageAmazon, "deploy storage amazon"))

	deployStorageGoogle := &cobra.Command{
		Use:    "{{alias}} <credentials-file>",
		Short:  "Deploy credentials for the Google Cloud storage provider.",
		Long:   "Deploy credentials for the Google Cloud storage provider, so that Pachyderm can ingress data from and egress data to it.",
		PreRun: preRun,
		Run: cmdutil.RunFixedArgs(1, func(args []string) error {
			credBytes, err := ioutil.ReadFile(args[0])
			if err != nil {
				return errors.Wrapf(err, "error reading credentials file %s", args[0])
			}
			return deployStorageSecrets(assets.GoogleSecret("", string(credBytes)))
		}),
	}
	appendGlobalFlags(deployStorageGoogle)
	commands = append(commands, cmdutil.CreateAlias(deployStorageGoogle, "deploy storage google"))

	deployStorageAzure := &cobra.Command{
		Use:    "{{alias}} <account-name> <account-key>",
		Short:  "Deploy credentials for the Azure storage provider.",
		Long:   "Deploy credentials for the Azure storage provider, so that Pachyderm can ingress data from and egress data to it.",
		PreRun: preRun,
		Run: cmdutil.RunFixedArgs(2, func(args []string) error {
			return deployStorageSecrets(assets.MicrosoftSecret("", args[0], args[1]))
		}),
	}
	appendGlobalFlags(deployStorageAzure)
	commands = append(commands, cmdutil.CreateAlias(deployStorageAzure, "deploy storage microsoft"))

	deployStorage := &cobra.Command{
		Short: "Deploy credentials for a particular storage provider.",
		Long:  "Deploy credentials for a particular storage provider, so that Pachyderm can ingress data from and egress data to it.",
	}
	commands = append(commands, cmdutil.CreateAlias(deployStorage, "deploy storage"))

	listImages := &cobra.Command{
		Short:  "Output the list of images in a deployment.",
		Long:   "Output the list of images in a deployment.",
		PreRun: preRun,
		Run: cmdutil.RunFixedArgs(0, func(args []string) error {
			for _, image := range assets.Images(opts) {
				fmt.Println(image)
			}
			return nil
		}),
	}
	appendGlobalFlags(listImages)
	commands = append(commands, cmdutil.CreateAlias(listImages, "deploy list-images"))

	exportImages := &cobra.Command{
		Use:    "{{alias}} <output-file>",
		Short:  "Export a tarball (to stdout) containing all of the images in a deployment.",
		Long:   "Export a tarball (to stdout) containing all of the images in a deployment.",
		PreRun: preRun,
		Run: cmdutil.RunFixedArgs(1, func(args []string) (retErr error) {
			file, err := os.Create(args[0])
			if err != nil {
				return err
			}
			defer func() {
				if err := file.Close(); err != nil && retErr == nil {
					retErr = err
				}
			}()
			return images.Export(opts, file)
		}),
	}
	appendGlobalFlags(exportImages)
	commands = append(commands, cmdutil.CreateAlias(exportImages, "deploy export-images"))

	importImages := &cobra.Command{
		Use:    "{{alias}} <input-file>",
		Short:  "Import a tarball (from stdin) containing all of the images in a deployment and push them to a private registry.",
		Long:   "Import a tarball (from stdin) containing all of the images in a deployment and push them to a private registry.",
		PreRun: preRun,
		Run: cmdutil.RunFixedArgs(1, func(args []string) (retErr error) {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer func() {
				if err := file.Close(); err != nil && retErr == nil {
					retErr = err
				}
			}()
			return images.Import(opts, file)
		}),
	}
	appendGlobalFlags(importImages)
	commands = append(commands, cmdutil.CreateAlias(importImages, "deploy import-images"))

	return commands
}

// Cmds returns a list of cobra commands for deploying Pachyderm clusters.
func Cmds() []*cobra.Command {
	commands := standardDeployCmds()

	var lbTLSHost string
	var lbTLSEmail string
	var dryRun bool
	var outputFormat string
	var hubImage string
	var userImage string
	deployIDE := &cobra.Command{
		Short: "Deploy the Pachyderm IDE.",
		Long:  "Deploy a JupyterHub-based IDE alongside the Pachyderm cluster.",
		Run: cmdutil.RunFixedArgs(0, func(args []string) (retErr error) {
			cfg, err := config.Read(false)
			if err != nil {
				return err
			}
			_, activeContext, err := cfg.ActiveContext(true)
			if err != nil {
				return err
			}

			c, err := client.NewOnUserMachine("user")
			if err != nil {
				return errors.Wrapf(err, "error constructing pachyderm client")
			}
			defer c.Close()

			enterpriseResp, err := c.Enterprise.GetState(c.Ctx(), &enterprise.GetStateRequest{})
			if err != nil {
				return errors.Wrapf(grpcutil.ScrubGRPC(err), "could not get Enterprise status")
			}

			if enterpriseResp.State != enterprise.State_ACTIVE {
				return errors.New("Pachyderm Enterprise must be enabled to use this feature")
			}

			authActive, err := c.IsAuthActive()
			if err != nil {
				return errors.Wrapf(grpcutil.ScrubGRPC(err), "could not check whether auth is active")
			}
			if !authActive {
				return errors.New("Pachyderm auth must be enabled to use this feature")
			}

			whoamiResp, err := c.WhoAmI(c.Ctx(), &auth.WhoAmIRequest{})
			if err != nil {
				return errors.Wrapf(grpcutil.ScrubGRPC(err), "could not get the current logged in user")
			}

			authTokenResp, err := c.GetAuthToken(c.Ctx(), &auth.GetAuthTokenRequest{
				Subject: whoamiResp.Username,
			})
			if err != nil {
				return errors.Wrapf(grpcutil.ScrubGRPC(err), "could not get an auth token")
			}

			hubImageName, hubImageTag := docker.ParseRepositoryTag(hubImage)
			userImageName, userImageTag := docker.ParseRepositoryTag(userImage)

			values := map[string]interface{}{
				"hub": map[string]interface{}{
					"image": map[string]interface{}{
						"name": hubImageName,
						"tag":  hubImageTag,
					},
					"extraConfig": map[string]interface{}{
						"jupyterlab": "c.Spawner.cmd = ['jupyter-labhub']",
						"templates":  "c.JupyterHub.template_paths = ['/app/templates']",
					},
				},
				"singleuser": map[string]interface{}{
					"image": map[string]interface{}{
						"name": userImageName,
						"tag":  userImageTag,
					},
					"defaultUrl": "/lab",
				},
				"auth": map[string]interface{}{
					"state": map[string]interface{}{
						"enabled":   true,
						"cryptoKey": generateSecureToken(16),
					},
					"type": "custom",
					"custom": map[string]interface{}{
						"className": "pachyderm_authenticator.PachydermAuthenticator",
						"config": map[string]interface{}{
							"pach_auth_token": authTokenResp.Token,
						},
					},
					"admin": map[string]interface{}{
						"users": []string{whoamiResp.Username},
					},
				},
				"proxy": map[string]interface{}{
					"secretToken": generateSecureToken(16),
				},
			}

			if lbTLSHost != "" && lbTLSEmail != "" {
				values["https"] = map[string]interface{}{
					"hosts": []string{lbTLSHost},
					"letsencrypt": map[string]interface{}{
						"contactEmail": lbTLSEmail,
					},
				}
			}

			if dryRun {
				var buf bytes.Buffer
				enc := encoder(outputFormat, &buf)
				if err = enc.Encode(values); err != nil {
					return err
				}
				_, err = os.Stdout.Write(buf.Bytes())
				return err
			}

			rel, err := helm.Deploy(
				activeContext,
				"jupyterhub",
				"https://jupyterhub.github.io/helm-chart/",
				"pachyderm-ide",
				"jupyterhub/jupyterhub",
				jupyterhubChartVersion,
				values,
			)
			if err != nil {
				return errors.Wrapf(err, "failed to deploy Pachyderm IDE")
			}

			fmt.Println(rel.Info.Notes)
			return nil
		}),
	}
	deployIDE.Flags().StringVar(&lbTLSHost, "lb-tls-host", "", "Hostname for minting a Let's Encrypt TLS cert on the load balancer")
	deployIDE.Flags().StringVar(&lbTLSEmail, "lb-tls-email", "", "Contact email for minting a Let's Encrypt TLS cert on the load balancer")
	deployIDE.Flags().BoolVar(&dryRun, "dry-run", false, "Don't actually deploy, instead just print the Helm config.")
	deployIDE.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format. One of: json|yaml")
	deployIDE.Flags().StringVar(&hubImage, "hub-image", defaultIDEHubImage, "Image for IDE hub")
	deployIDE.Flags().StringVar(&userImage, "user-image", defaultIDEUserImage, "Image for IDE user environments")
	commands = append(commands, cmdutil.CreateAlias(deployIDE, "deploy ide"))

	deploy := &cobra.Command{
		Short: "Deploy a Pachyderm cluster.",
		Long:  "Deploy a Pachyderm cluster.",
	}
	commands = append(commands, cmdutil.CreateAlias(deploy, "deploy"))

	var all bool
	var includingMetadata bool
	var includingIDE bool
	var namespace string
	undeploy := &cobra.Command{
		Short: "Tear down a deployed Pachyderm cluster.",
		Long:  "Tear down a deployed Pachyderm cluster.",
		Run: cmdutil.RunFixedArgs(0, func(args []string) error {
			// TODO(ys): remove the `--namespace` flag here eventually
			if namespace != "" {
				fmt.Printf("WARNING: The `--namespace` flag is deprecated and will be removed in a future version. Please set the namespace in the pachyderm context instead: pachctl config update context `pachctl config get active-context` --namespace '%s'\n", namespace)
			}
			// TODO(ys): remove the `--all` flag here eventually
			if all {
				fmt.Printf("WARNING: The `--all` flag is deprecated and will be removed in a future version. Please use `--metadata` instead.\n")
				includingMetadata = true
			}

			if includingMetadata {
				fmt.Printf(`
You are going to delete persistent volumes where metadata is stored. If your
persistent volumes were dynamically provisioned (i.e. if you used the
"--dynamic-etcd-nodes" flag), the underlying volumes will be removed, making
metadata such as repos, commits, pipelines, and jobs unrecoverable. If your
persistent volume was manually provisioned (i.e. if you used the
"--static-etcd-volume" flag), the underlying volume will not be removed.
`)
			}

			if ok, err := cmdutil.InteractiveConfirm(); err != nil {
				return err
			} else if !ok {
				return nil
			}

			cfg, err := config.Read(false)
			if err != nil {
				return err
			}
			_, activeContext, err := cfg.ActiveContext(true)
			if err != nil {
				return err
			}

			if namespace == "" {
				namespace = activeContext.Namespace
			}

			assets := []string{
				"service",
				"replicationcontroller",
				"deployment",
				"serviceaccount",
				"secret",
				"statefulset",
				"clusterrole",
				"clusterrolebinding",
			}
			if includingMetadata {
				assets = append(assets, []string{
					"storageclass",
					"persistentvolumeclaim",
					"persistentvolume",
				}...)
			}
			if err := kubectl(nil, activeContext, "delete", strings.Join(assets, ","), "-l", "suite=pachyderm", "--namespace", namespace); err != nil {
				return err
			}

			if includingIDE {
				// remove IDE
				if err = helm.Destroy(activeContext, "pachyderm-ide", namespace); err != nil {
					log.Errorf("failed to delete helm installation: %v", err)
				}
				ideAssets := []string{
					"replicaset",
					"deployment",
					"service",
					"pod",
				}
				if err = kubectl(nil, activeContext, "delete", strings.Join(ideAssets, ","), "-l", "app=jupyterhub", "--namespace", namespace); err != nil {
					return err
				}
			}

			// remove the context from the config
			kubeConfig, err := config.RawKubeConfig()
			if err != nil {
				return err
			}
			kubeContext := kubeConfig.Contexts[kubeConfig.CurrentContext]
			if kubeContext != nil {
				cfg, err := config.Read(true)
				if err != nil {
					return err
				}
				ctx := &config.Context{
					ClusterName: kubeContext.Cluster,
					AuthInfo:    kubeContext.AuthInfo,
					Namespace:   namespace,
				}

				// remove _all_ contexts associated with this
				// deployment
				configUpdated := false
				for {
					contextName, _ := findEquivalentContext(cfg, ctx)
					if contextName == "" {
						break
					}
					configUpdated = true
					delete(cfg.V2.Contexts, contextName)
					if contextName == cfg.V2.ActiveContext {
						cfg.V2.ActiveContext = ""
					}
				}
				if configUpdated {
					if err = cfg.Write(); err != nil {
						return err
					}
				}
			}

			return nil
		}),
	}
	undeploy.Flags().BoolVarP(&all, "all", "a", false, "DEPRECATED: Use \"--metadata\" instead.")
	undeploy.Flags().BoolVarP(&includingMetadata, "metadata", "", false, `
Delete persistent volumes where metadata is stored. If your persistent volumes
were dynamically provisioned (i.e. if you used the "--dynamic-etcd-nodes"
flag), the underlying volumes will be removed, making metadata such as repos,
commits, pipelines, and jobs unrecoverable. If your persistent volume was
manually provisioned (i.e. if you used the "--static-etcd-volume" flag), the
underlying volume will not be removed.`)
	undeploy.Flags().BoolVarP(&includingIDE, "ide", "", false, "Delete the Pachyderm IDE deployment if it exists.")
	undeploy.Flags().StringVar(&namespace, "namespace", "", "Kubernetes namespace to undeploy Pachyderm from.")
	commands = append(commands, cmdutil.CreateAlias(undeploy, "undeploy"))

	var updateDashDryRun bool
	var updateDashOutputFormat string
	updateDash := &cobra.Command{
		Short: "Update and redeploy the Pachyderm Dashboard at the latest compatible version.",
		Long:  "Update and redeploy the Pachyderm Dashboard at the latest compatible version.",
		Run: cmdutil.RunFixedArgs(0, func(args []string) error {
			cfg, err := config.Read(false)
			if err != nil {
				return err
			}
			_, activeContext, err := cfg.ActiveContext(false)
			if err != nil {
				return err
			}

			// Undeploy the dash
			if !updateDashDryRun {
				if err := kubectl(nil, activeContext, "delete", "deploy", "-l", "suite=pachyderm,app=dash"); err != nil {
					return err
				}
				if err := kubectl(nil, activeContext, "delete", "svc", "-l", "suite=pachyderm,app=dash"); err != nil {
					return err
				}
			}
			// Redeploy the dash
			var buf bytes.Buffer
			opts := &assets.AssetOpts{
				DashOnly:  true,
				DashImage: getDefaultOrLatestDashImage("", updateDashDryRun),
			}
			if err := assets.WriteDashboardAssets(
				encoder(updateDashOutputFormat, &buf), opts,
			); err != nil {
				return err
			}
			return kubectlCreate(updateDashDryRun, buf.Bytes(), opts)
		}),
	}
	updateDash.Flags().BoolVar(&updateDashDryRun, "dry-run", false, "Don't actually deploy Pachyderm Dash to Kubernetes, instead just print the manifest.")
	updateDash.Flags().StringVarP(&updateDashOutputFormat, "output", "o", "json", "Output format. One of: json|yaml")
	commands = append(commands, cmdutil.CreateAlias(updateDash, "update-dash"))

	return commands
}

func getDefaultOrLatestDashImage(dashImage string, dryRun bool) string {
	var err error
	version := version.PrettyPrintVersion(version.Version)
	defer func() {
		if err != nil && !dryRun {
			fmt.Printf("No updated dash image found for pachctl %v: %v Falling back to dash image %v\n", version, err, defaultDashImage)
		}
	}()
	if dashImage != "" {
		// It has been supplied explicitly by version on the command line
		return dashImage
	}
	dashImage = defaultDashImage
	compatibleDashVersionsURL := fmt.Sprintf("https://raw.githubusercontent.com/pachyderm/pachyderm/master/etc/compatibility/%v", version)
	resp, err := http.Get(compatibleDashVersionsURL)
	if err != nil {
		return dashImage
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dashImage
	}
	if resp.StatusCode != 200 {
		err = errors.New(string(body))
		return dashImage
	}
	allVersions := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(allVersions) < 1 {
		return dashImage
	}
	latestVersion := strings.TrimSpace(allVersions[len(allVersions)-1])

	return fmt.Sprintf("pachyderm/dash:%v", latestVersion)
}
