## SeujtaCita Tech test

Rest API CRUD User dan User Login, dengan JWT token/refreshToken

### App Stacks
 - Backend: golang dengan framework fiber, mysql dengan orm gorm
 - Deployments: docker, kubernetes. Kubernetes jalan pada pada ubuntu desktop dengan installer dari minikube.

#### Menjalankan aplikasi dengan docker

 - dev
	```
	$ docker-compose up -d --build
	```
 - prod
	```
	$ docker build -t sejuta-cita .
	$ docker run -d --name sejuta-cita -e <ENV> -p 3000:3000 sejuta-cita 
	```
	envoiremnts silahkan lihat di file docker-compose.yml

#### Menjalankan aplikasi dengan kubernetes cluster
 1. **setup mysql** 
 	```
	$ kubectl create -f k8s.mysql.secret.yaml
	$ kubectl apply -f k8s.mysql.pv.yaml
	$ kubectl apply -f k8s.mysql.pvc.yaml
	$ kubectl apply -f k8s.mysql.deployment.yml
	$ kubectl apply -f k8s.mysql.service.yml
	```
 2. **Setup app**
	aplikasi telah di build dan tersimpan di image registery docker hub, dengan nama 
	emixbal/sejuta-cita gunakan tag v1
	``` emixbal/sejuta-cita:v1 ```
 	```
	$ kubectl apply -f k8s.app.deployment.yml
	$ kubectl apply -f k8s.app.service.yml
	```
 3. **Check services dan pods apakah sukses?**
    - Menampilkan semua services
        ```
        $ kubectl get services
        ```
        ```
        $ kubectl get pods
        ```  
    - karena menggunkan minikube harus set external ip secarea manual
        ```
        $ kubectl service sejuta-cita-service --url
        ```
        atau menngunakan minikube tunnel
 6. **Set url yg muncul sebagai baseUrl**

