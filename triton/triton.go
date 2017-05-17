package triton

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/golang/glog"
	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"

	"k8s.io/client-go/tools/cache"

	api "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"k8s.io/apimachinery/pkg/types"
)

type Host struct {
	TritonHost  *triton.Machine
	IPAddresses []string
}

type PublicEndpoint struct {
	IPAddress string
	Port      int
}

const (
	providerName = "triton"
)

// CloudProvider is an interface to the Triton Go SDK (`triton-go`)
type CloudProvider struct {
	client    *triton.Client
	conf      *tConfig
	hostCache cache.Store
}

// ProviderName returns the cloud provider ID.
func (t *CloudProvider) ProviderName() string {
	return providerName
}

// Instances return an implementation of Instances for Triton
func (t *CloudProvider) Instances() (cloudprovider.Instances, bool) {
	return t, true
}

// --- Instances Functions ---

// NodeAddresses returns the addresses of the specified instance.
//
// This implementation only returns the address of the calling instance. This is
// ok because the gce implementation makes that assumption and the comment for
// the interface states it as a todo to clarify that it is only for the current
// host
func (t *CloudProvider) NodeAddresses(name types.NodeName) ([]api.NodeAddress, error) {
	host, err := t.hostGetOrFetchFromCache(string(name))
	if err != nil {
		return nil, err
	}

	addresses := []api.NodeAddress{}
	for _, ip := range host.IPAddresses {
		addresses = append(addresses, api.NodeAddress{Type: api.NodeExternalIP, Address: ip.Address})
		addresses = append(addresses, api.NodeAddress{Type: api.NodeLegacyHostIP, Address: ip.Address})
	}
	addresses = append(addresses, api.NodeAddress{Type: api.NodeHostName, Address: host.TritonHost.Hostname})

	return addresses, nil
}

// ExternalID returns the cloud provider ID of the specified instance
// (deprecated).
func (t *CloudProvider) ExternalID(name types.NodeName) (string, error) {
	glog.Infof("ExternalID [%s]", string(name))
	return t.InstanceID(name)
}

// InstanceID returns the cloud provider ID of the specified instance.
func (t *CloudProvider) InstanceID(name types.NodeName) (string, error) {
	glog.Infof("InstanceID [%s]", string(name))
	host, err := t.hostGetOrFetchFromCache(string(name))
	if err != nil {
		return "", err
	}

	return host.TritonHost.Uuid, nil
}

// InstanceType returns the type of the specified instance.
//
// Note that if the instance does not exist or is no longer running, we must
// return ("", cloudprovider.InstanceNotFound)
func (t *CloudProvider) InstanceType(name types.NodeName) (string, error) {
	_, err := t.InstanceID(name)
	if err != nil {
		return "", err
	}

	// TODO: Triton Machine types have a Brand attribute, we can return that here
	return "triton", nil
}

// List lists instances that match 'filter' which is a regular expression which
// must match the entire instance name (fqdn)
func (t *CloudProvider) List(filter string) ([]types.NodeName, error) {
	glog.Infof("List %s", filter)

	// opts := client.NewListOpts()
	// opts.Filters["removed_null"] = "1"

	// hosts, err := t.client.Host.List(opts)
	// if err != nil {
	//  return nil, fmt.Errorf("Coudln't get hosts by filter [%s]. Error: %#v", filter, err)
	// }

	// if len(hosts.Data) == 0 {
	// 	return nil, fmt.Errorf("No hosts found")
	// }

	// if strings.HasPrefix(filter, "'") && strings.HasSuffix(filter, "'") {
	// 	filter = filter[1 : len(filter)-1]
	// }

	// re, err := regexp.Compile(filter)
	// if err != nil {
	// 	return nil, err
	// }

	retHosts := []types.NodeName{}
	// for _, host := range hosts.Data {
	// 	if re.MatchString(host.Hostname) {
	// 		retHosts = append(retHosts, types.NodeName(host.Hostname))
	// 	}
	// }

	return retHosts, err
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
//
// TODO: This might be doable considering Triton reuses a user's private SSH key
// and can derive a public key then upload to each k8s host node.
func (t *CloudProvider) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return fmt.Errorf("Not implemented")
}

// CurrentNodeName returns the name of the node we are currently running on
func (t *CloudProvider) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (t *CloudProvider) addHostToCache(host *Host) {
	if host != nil {
		t.hostCache.Add(host)
	}
}

func (t *CloudProvider) removeFromCache(name string) {
	host := t.getHostFromCache(name)
	if host != nil {
		t.hostCache.Delete(host)
	}
}

func (t *CloudProvider) getHostFromCache(name string) *Host {
	var host *Host

	// entry gets expired once retrieved
	defer t.addHostToCache(host)

	hostObj, exists, err := t.hostCache.GetByKey(name)
	if err == nil && exists {
		h, ok := hostObj.(*Host)
		if ok {
			host = h
		}
	}
	return host
}

func (t *CloudProvider) hostGetOrFetchFromCache(name string) (*Host, error) {
	host, err := t.getHostByName(name)
	if err != nil {
		if err == cloudprovider.InstanceNotFound {
			// evict from cache
			t.removeFromCache(name)
			return nil, err
		} else {
			host := t.getHostFromCache(name)
			if host != nil {
				return host, nil
			}
		}
	}
	t.addHostToCache(host)
	return host, nil
}

func (t *CloudProvider) getHostByName(name string) (*Host, error) {
	// opts := client.NewListOpts()
	// opts.Filters["removed_null"] = "1"
	// hosts, err := t.client.Host.List(opts)
	// if err != nil {
	// 	return nil, fmt.Errorf("Coudln't get host by name [%s]. Error: %#v", name, err)
	// }

	// hostsToReturn := make([]client.Host, 0)
	// for _, host := range hosts.Data {
	// 	if strings.EqualFold(host.Hostname, name) {
	// 		hostsToReturn = append(hostsToReturn, host)
	// 	}
	// }

	// if len(hostsToReturn) == 0 {
	// 	return nil, cloudprovider.InstanceNotFound
	// }

	// if len(hostsToReturn) > 1 {
	// 	return nil, fmt.Errorf("multiple instances found for name: %s", name)
	// }

	// tritonHost := &hostsToReturn[0]

	// coll := &client.IpAddressCollection{}
	// err = t.client.GetLink(tritonHost.Resource, "ipAddresses", coll)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error getting ip addresses for node [%s]. Error: %#v", name, err)
	// }

	// if len(coll.Data) == 0 {
	// 	return nil, cloudprovider.InstanceNotFound
	// }

	// host := &Host{
	// 	TritonHost:  tritonHost,
	// 	IPAddresses: coll.Data,
	// }

	return &Host{}, nil
}

// --- Utility functions ---

func Init(configFilePath string) (cloudprovider.Interface, error) {
	if configFilePath != "" {
		var config *os.File
		config, err := os.Open(configFilePath)
		if err != nil {
			glog.Fatalf("Couldn't open cloud provider configuration %s: %#v",
				configFilePath, err)
		}

		defer config.Close()
		return newTritonCloud(config)
	}
	return newTritonCloud(nil)
}

type configGlobal struct {
	TritonEndpoint   string `gcfg:"triton-endpoint"`
	TritonAccount    string `gcfg:"triton-account"`
	TritonKeyId      string `gcfg:"triton-key-id"`
	TritonPrivateKey string `gcfg:"triton-private-key"`
}

type tConfig struct {
	Global configGlobal
}

func newTritonCloud(config io.Reader) (cloudprovider.Interface, error) {
	endpoint := os.Getenv("SDC_ENDPOINT")
	account := os.Getenv("SDC_ACCOUNT")
	keyID := os.Getenv("SDC_KEY_ID")
	privateKey := os.Getenv("SDC_PRIVATE_KEY")
	conf := tConfig{
		Global: configGlobal{
			TritonEndpoint:   endpoint,
			TritonAccount:    account,
			TritonKeyId:      keyID,
			TritonPrivateKey: privateKey,
		},
	}

	client, err := getTritonClient(conf)
	if err != nil {
		return nil, fmt.Errorf("Could not create Triton client: %#v", err)
	}

	cache := cache.NewTTLStore(hostStoreKeyFunc, time.Duration(24)*time.Hour)

	return &CloudProvider{
		client:    client,
		conf:      &conf,
		hostCache: cache,
	}, nil
}

func hostStoreKeyFunc(obj interface{}) (string, error) {
	return obj.(*Host).TritonHost.Name, nil
}

func getTritonClient(conf tConfig) (*triton.Client, error) {
	keyID := conf.Global.TritonKeyId
	accountName := conf.Global.TritonAccount
	privateKey := conf.Global.TritonPrivateKey

	sshKeySigner, err := authentication.NewPrivateKeySigner(keyID, privateKey, accountName)
	if err != nil {
		log.Fatalf("Fatal exception from NewPrivateKeySigner: %s", err)
	}

	newTritonClient := triton.NewClient(conf.Global.TritonEndpoint, conf.Global.TritonAccount, sshKeySigner)
	return newTritonClient
}
