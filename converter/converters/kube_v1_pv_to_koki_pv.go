package converters

import (
	"reflect"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Kube_v1_PersistentVolume_to_Koki_PersistentVolume(kubePV *v1.PersistentVolume) (*types.PersistentVolumeWrapper, error) {
	var err error
	kokiPV := &types.PersistentVolume{}

	kokiPV.Name = kubePV.Name
	kokiPV.Namespace = kubePV.Namespace
	kokiPV.Version = kubePV.APIVersion
	kokiPV.Cluster = kubePV.ClusterName
	kokiPV.Labels = kubePV.Labels
	kokiPV.Annotations = kubePV.Annotations

	kubeSpec := kubePV.Spec
	kokiPV.Storage, err = convertCapacity(kubeSpec.Capacity)
	if err != nil {
		return nil, err
	}

	kokiPV.PersistentVolumeSource, err = convertPersistentVolumeSource(kubeSpec.PersistentVolumeSource)
	if err != nil {
		return nil, err
	}
	if len(kubeSpec.AccessModes) > 0 {
		kokiPV.AccessModes = &types.AccessModes{
			Modes: kubeSpec.AccessModes,
		}
	}
	kokiPV.Claim = kubeSpec.ClaimRef
	kokiPV.ReclaimPolicy = convertReclaimPolicy(kubeSpec.PersistentVolumeReclaimPolicy)
	kokiPV.StorageClass = kubeSpec.StorageClassName
	if len(kubeSpec.MountOptions) > 0 {
		kokiPV.MountOptions = strings.Join(kubeSpec.MountOptions, ",")
	}

	if !reflect.DeepEqual(kubePV.Status, v1.PersistentVolumeStatus{}) {
		kokiPV.Status = &kubePV.Status
	}

	return &types.PersistentVolumeWrapper{
		PersistentVolume: *kokiPV,
	}, nil
}

func convertSecretReference(kubeRef *v1.SecretReference) *types.SecretReference {
	if kubeRef == nil {
		return nil
	}

	return &types.SecretReference{
		Namespace: kubeRef.Namespace,
		Name:      kubeRef.Name,
	}
}

func convertCephFSPersistentSecretFileOrRef(kubeFile string, kubeRef *v1.SecretReference) *types.CephFSPersistentSecretFileOrRef {
	if len(kubeFile) > 0 {
		return &types.CephFSPersistentSecretFileOrRef{
			File: kubeFile,
		}
	}

	if kubeRef != nil {
		return &types.CephFSPersistentSecretFileOrRef{
			Ref: convertSecretReference(kubeRef),
		}
	}

	return nil
}

func convertPersistentVolumeSource(kubeSource v1.PersistentVolumeSource) (types.PersistentVolumeSource, error) {
	if kubeSource.GCEPersistentDisk != nil {
		return types.PersistentVolumeSource{
			GcePD: convertGcePDVolume(kubeSource.GCEPersistentDisk),
		}, nil
	}
	if kubeSource.AWSElasticBlockStore != nil {
		return types.PersistentVolumeSource{
			AwsEBS: convertAwsEBSVolume(kubeSource.AWSElasticBlockStore),
		}, nil
	}
	if kubeSource.HostPath != nil {
		source, err := convertHostPathVolume(kubeSource.HostPath)
		if err != nil {
			return types.PersistentVolumeSource{}, err
		}
		return types.PersistentVolumeSource{
			HostPath: source,
		}, nil
	}
	if kubeSource.Glusterfs != nil {
		return types.PersistentVolumeSource{
			Glusterfs: convertGlusterfsVolume(kubeSource.Glusterfs),
		}, nil
	}
	if kubeSource.NFS != nil {
		return types.PersistentVolumeSource{
			NFS: convertNFSVolume(kubeSource.NFS),
		}, nil
	}
	if kubeSource.ISCSI != nil {
		return types.PersistentVolumeSource{
			ISCSI: convertISCSIVolume(kubeSource.ISCSI),
		}, nil
	}
	if kubeSource.Cinder != nil {
		return types.PersistentVolumeSource{
			Cinder: convertCinderVolume(kubeSource.Cinder),
		}, nil
	}
	if kubeSource.FC != nil {
		return types.PersistentVolumeSource{
			FibreChannel: convertFibreChannelVolume(kubeSource.FC),
		}, nil
	}
	if kubeSource.Flocker != nil {
		return types.PersistentVolumeSource{
			Flocker: convertFlockerVolume(kubeSource.Flocker),
		}, nil
	}
	if kubeSource.FlexVolume != nil {
		return types.PersistentVolumeSource{
			Flex: convertFlexVolume(kubeSource.FlexVolume),
		}, nil
	}
	if kubeSource.VsphereVolume != nil {
		return types.PersistentVolumeSource{
			Vsphere: convertVsphereVolume(kubeSource.VsphereVolume),
		}, nil
	}
	if kubeSource.Quobyte != nil {
		return types.PersistentVolumeSource{
			Quobyte: convertQuobyteVolume(kubeSource.Quobyte),
		}, nil
	}
	if kubeSource.AzureDisk != nil {
		source, err := convertAzureDiskVolume(kubeSource.AzureDisk)
		if err != nil {
			return types.PersistentVolumeSource{}, err
		}
		return types.PersistentVolumeSource{
			AzureDisk: source,
		}, nil
	}
	if kubeSource.PhotonPersistentDisk != nil {
		return types.PersistentVolumeSource{
			PhotonPD: convertPhotonPDVolume(kubeSource.PhotonPersistentDisk),
		}, nil
	}
	if kubeSource.PortworxVolume != nil {
		return types.PersistentVolumeSource{
			Portworx: convertPortworxVolume(kubeSource.PortworxVolume),
		}, nil
	}
	if kubeSource.RBD != nil {
		source := kubeSource.RBD
		return types.PersistentVolumeSource{
			RBD: &types.RBDPersistentVolume{
				CephMonitors: source.CephMonitors,
				RBDImage:     source.RBDImage,
				FSType:       source.FSType,
				RBDPool:      source.RBDPool,
				RadosUser:    source.RadosUser,
				Keyring:      source.Keyring,
				SecretRef:    convertSecretReference(source.SecretRef),
				ReadOnly:     source.ReadOnly,
			},
		}, nil
	}
	if kubeSource.CephFS != nil {
		source := kubeSource.CephFS
		secretFileOrRef := convertCephFSPersistentSecretFileOrRef(source.SecretFile, source.SecretRef)
		return types.PersistentVolumeSource{
			CephFS: &types.CephFSPersistentVolume{
				Monitors:        source.Monitors,
				Path:            source.Path,
				User:            source.User,
				SecretFileOrRef: secretFileOrRef,
				ReadOnly:        source.ReadOnly,
			},
		}, nil
	}
	if kubeSource.AzureFile != nil {
		source := kubeSource.AzureFile
		return types.PersistentVolumeSource{
			AzureFile: &types.AzureFilePersistentVolume{
				Secret: types.SecretReference{
					Name:      source.SecretName,
					Namespace: util.FromStringPtr(source.SecretNamespace),
				},
				ShareName: source.ShareName,
				ReadOnly:  source.ReadOnly,
			},
		}, nil
	}

	return types.PersistentVolumeSource{}, util.InvalidInstanceErrorf(kubeSource, "didn't find any supported volume source")
}

func convertReclaimPolicy(kubePolicy v1.PersistentVolumeReclaimPolicy) types.PersistentVolumeReclaimPolicy {
	return types.PersistentVolumeReclaimPolicy(strings.ToLower(string(kubePolicy)))
}

func convertCapacity(kubeCapacity v1.ResourceList) (*resource.Quantity, error) {
	if len(kubeCapacity) == 0 {
		return nil, nil
	}

	for res, quantity := range kubeCapacity {
		if res == v1.ResourceStorage {
			return &quantity, nil
		}
	}

	return nil, util.InvalidInstanceErrorf(kubeCapacity, "only supports Storage resource")
}
