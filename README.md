# Install Gloo Portal

## Add the helm repo
```
helm repo add dev-portal https://storage.googleapis.com/dev-portal-helm
helm repo update
```

## Create helm values override
```
cat << EOF > gloo-values.yaml
gloo:
  enabled: true # Enables integration with Gloo Edge Enterprise
licenseKey:
  secretRef:
    name: license
    namespace: gloo-system
    key: license-key
EOF
```

## Create the namespace and install the helm chart
```
k create ns dev-portal
helm install dev-portal dev-portal/dev-portal -n dev-portal --values gloo-values.yaml
```