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
        ![all service](https://raw.githubusercontent.com/emixbal/sejuta-cita/main/images/services%20all.png)
        
        ```
        $ kubectl get pods
        ```  
        ![all pods](https://raw.githubusercontent.com/emixbal/sejuta-cita/main/images/pods%20all.png)
    - karena menggunkan minikube harus set external ip secarea manual
        ```
        $ kubectl service sejuta-cita-service --url
        ```
        ![all pods](https://raw.githubusercontent.com/emixbal/sejuta-cita/main/images/services%20generate%20url.png)
        atau menngunakan minikube tunnel
 6. **Set url yg muncul sebagai baseUrl**

#### Mengakses API dengan Postman client
 1. download postman colections dari link ini
    [https://www.getpostman.com/collections/3797c3347deb99272049](https://www.getpostman.com/collections/3797c3347deb99272049)
 2. saat app dijalankan telah otomati dibuat user level admin seeder. dengan credential
    - email=emixbal@gmail.com
    - password=aaaaaaaa
    dengan credential diatas gunakan request "login refresh token" untuk login
 3. ketika mendapat access access_token & refresh_token buat envoirement, dengan key
    - baseUrl, lalau isikan value link yg telah digenarate
    - jwtToken, lalu isikan dengan access_token yg didapat
    - jwtRefreshToken, lalu isikan dengan refresh_token yg didapat

