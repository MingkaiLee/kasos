git submodule update --init
cd ./kube-prometheus
minikube start --kubernetes-version=v1.28.0 --memory=12g --cpus=8 --bootstrapper=kubeadm --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.bind-address=0.0.0.0 --extra-config=controller-manager.bind-address=0.0.0.0
cd ..
cd ./build
./build.sh