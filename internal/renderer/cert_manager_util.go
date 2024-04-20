package renderer

import (
	acmev1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	cert "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var baseSecret = corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "",
		Namespace: "",
		Labels: map[string]string{
			opdefault.OwnedByLabelKey: opdefault.OwnedByLabelValue,
			//opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
			opdefault.RelatedComponent: opdefault.ComponentApplicationExpose,
		},
	},
	Data: map[string][]byte{},
	Type: corev1.SecretTypeOpaque,
}

var baseIssuer = cert.Issuer{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "",
		Namespace: "",
		Labels: map[string]string{
			opdefault.OwnedByLabelKey: opdefault.OwnedByLabelValue,
			//opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
			opdefault.RelatedComponent: opdefault.ComponentApplicationExpose,
		},
	},
	Spec: cert.IssuerSpec{IssuerConfig: cert.IssuerConfig{
		ACME: &acmev1.ACMEIssuer{
			//TODO fix email
			Email:          "replace@withyourown.io",
			Server:         "https://acme-v02.api.letsencrypt.org/directory",
			PreferredChain: "",
			PrivateKey: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "cloudflare-issuer-account-key",
				},
			},
			Solvers: []acmev1.ACMEChallengeSolver{{
				Selector: nil,
				DNS01:    nil,
			}},
		},
	}},
}

func createIssuer(lkMesh lkstnv1a1.LiveKitMesh) (*cert.Issuer, *corev1.Secret) {
	appExpose := lkMesh.Spec.Components.ApplicationExpose
	switch *appExpose.CertManager.Issuer.ChallengeSolver {
	//TODO create func for each issuertype
	case opdefault.IssuerCloudFlare:
		return newCloudFlareIssuer(appExpose, lkMesh.GetNamespace())
	case opdefault.IssuerRoute53:
		return nil, nil
	case opdefault.IssuerAzureDNS:
		return nil, nil
	case opdefault.IssuerDigitalOcean:
		return nil, nil
	case opdefault.IssuerCloudDNS:
		return nil, nil
	default:
		panic("invalid baseIssuer type, if you got this the API has a bug")
	}
}

func newCloudFlareIssuer(appExpose *lkstnv1a1.ApplicationExpose, namespace string) (*cert.Issuer, *corev1.Secret) {
	//setting Issuer
	certManager := appExpose.CertManager
	issuer := baseIssuer
	issuer.Name = "cloudflare-issuer"
	issuer.Namespace = namespace
	issuer.Labels[opdefault.RelatedLiveKitMeshKey] = opdefault.ComponentApplicationExpose

	solver := &issuer.Spec.IssuerConfig.ACME.Solvers[0]
	solver.DNS01 = &acmev1.ACMEChallengeSolverDNS01{
		Cloudflare: &acmev1.ACMEIssuerDNS01ProviderCloudflare{
			APIToken: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "cloudflare-api-token-secret",
				},
				Key: opdefault.DefaultClusterIssuerSecretApiTokenKey,
			},
		},
	}
	solver.Selector = &acmev1.CertificateDNSNameSelector{
		DNSZones: []string{
			*appExpose.HostName,
		},
	}

	if certManager.Issuer.Email != nil {
		issuer.Spec.IssuerConfig.ACME.Email = *certManager.Issuer.Email
	}

	// setting Secret
	secret := baseSecret
	secret.Name = "cloudflare-api-token-secret"
	secret.Namespace = namespace
	secret.Labels[opdefault.RelatedLiveKitMeshKey] = opdefault.ComponentApplicationExpose

	secret.Data[opdefault.DefaultClusterIssuerSecretApiTokenKey] = []byte(*certManager.Issuer.ApiToken)

	return &issuer, &secret
}
