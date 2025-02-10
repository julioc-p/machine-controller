/*
Copyright 2019 The Machine Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudprovider

import (
	"errors"

	cloudprovidercache "k8c.io/machine-controller/pkg/cloudprovider/cache"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/alibaba"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/anexia"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/aws"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/azure"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/baremetal"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/digitalocean"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/edge"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/equinixmetal"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/external"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/fake"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/gce"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/hetzner"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/kubevirt"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/linode"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/nutanix"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/opennebula"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/openstack"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/scaleway"
	vcd "k8c.io/machine-controller/pkg/cloudprovider/provider/vmwareclouddirector"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/vsphere"
	"k8c.io/machine-controller/pkg/cloudprovider/provider/vultr"
	cloudprovidertypes "k8c.io/machine-controller/pkg/cloudprovider/types"
	"k8c.io/machine-controller/pkg/providerconfig"
	providerconfigtypes "k8c.io/machine-controller/pkg/providerconfig/types"
)

var (
	cache = cloudprovidercache.New()

	// ErrProviderNotFound tells that the requested cloud provider was not found.
	ErrProviderNotFound = errors.New("cloudprovider not found")

	providers = map[providerconfigtypes.CloudProvider]func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider{
		providerconfigtypes.CloudProviderDigitalocean: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return digitalocean.New(cvr)
		},
		providerconfigtypes.CloudProviderAWS: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return aws.New(cvr)
		},
		providerconfigtypes.CloudProviderOpenstack: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return openstack.New(cvr)
		},
		providerconfigtypes.CloudProviderGoogle: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return gce.New(cvr)
		},
		providerconfigtypes.CloudProviderHetzner: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return hetzner.New(cvr)
		},
		providerconfigtypes.CloudProviderVsphere: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return vsphere.New(cvr)
		},
		providerconfigtypes.CloudProviderAzure: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return azure.New(cvr)
		},
		providerconfigtypes.CloudProviderEquinixMetal: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return equinixmetal.New(cvr)
		},
		// NB: This is explicitly left to allow old Packet machines to be deleted.
		// We can handle those machines in the same way as Equinix Metal machines
		// because there are no API changes.
		// TODO: Remove this after deprecation period.
		providerconfigtypes.CloudProviderPacket: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return equinixmetal.New(cvr)
		},
		providerconfigtypes.CloudProviderFake: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return fake.New(cvr)
		},
		providerconfigtypes.CloudProviderEdge: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return edge.New(cvr)
		},
		providerconfigtypes.CloudProviderKubeVirt: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return kubevirt.New(cvr)
		},
		providerconfigtypes.CloudProviderAlibaba: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return alibaba.New(cvr)
		},
		providerconfigtypes.CloudProviderScaleway: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return scaleway.New(cvr)
		},
		providerconfigtypes.CloudProviderAnexia: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return anexia.New(cvr)
		},
		providerconfigtypes.CloudProviderBaremetal: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			// TODO(MQ): add a baremetal driver.
			return baremetal.New(cvr)
		},
		providerconfigtypes.CloudProviderNutanix: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return nutanix.New(cvr)
		},
		providerconfigtypes.CloudProviderVMwareCloudDirector: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return vcd.New(cvr)
		},
		providerconfigtypes.CloudProviderExternal: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return external.New(cvr)
		},
	}

	// communityProviders holds a map of cloud providers that have been implemented by community members and
	// contributed to machine-controller. They are not end-to-end tested by the machine-controller development team.
	communityProviders = map[providerconfigtypes.CloudProvider]func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider{
		providerconfigtypes.CloudProviderLinode: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return linode.New(cvr)
		},
		providerconfigtypes.CloudProviderVultr: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return vultr.New(cvr)
		},
		providerconfigtypes.CloudProviderOpenNebula: func(cvr *providerconfig.ConfigVarResolver) cloudprovidertypes.Provider {
			return opennebula.New(cvr)
		},
	}
)

// ForProvider returns a CloudProvider actuator for the requested provider.
func ForProvider(p providerconfigtypes.CloudProvider, cvr *providerconfig.ConfigVarResolver) (cloudprovidertypes.Provider, error) {
	if p, found := providers[p]; found {
		return NewValidationCacheWrappingCloudProvider(p(cvr)), nil
	}
	if p, found := communityProviders[p]; found {
		return NewValidationCacheWrappingCloudProvider(p(cvr)), nil
	}
	return nil, ErrProviderNotFound
}
