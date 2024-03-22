# Sistemas Operativos 1
##  Proyecto 1 - Máquina Virtual de Monitoreo

![Javascript](https://img.shields.io/badge/-Typescript-0d0d0d?style=flat&logo=typescript)![Node.js](https://img.shields.io/badge/-Node.js-0d0d0d?style=flat&logo=node.js&logoColor=0CFE00)![React](https://img.shields.io/badge/-React-0d0d0d?style=flat&logo=react)![Vite](https://img.shields.io/badge/-Vite-0d0d0d?style=flat&logo=vite)![CSS](https://img.shields.io/badge/-CSS-0d0d0d?style=flat&logo=CSS3&logoColor=1575FF)![C](https://img.shields.io/badge/-C-0d0d0d?style=flat&logo=C&logoColor=1575FF)![Go](https://img.shields.io/badge/-Golang-0d0d0d?style=flat&logo=Go)
![Ubuntu](https://img.shields.io/badge/-Ubuntu-0d0d0d?style=flat&logo=Ubuntu)![Git](https://img.shields.io/badge/-Git-0d0d0d?style=flat&logo=Git)![Github](https://img.shields.io/badge/-Github-0d0d0d?style=flat&logo=Github)![Vscode](https://img.shields.io/badge/-VSCode-0d0d0d?style=flat&logo=VisualStudioCode&logoColor=1575FF)![VMWare](https://img.shields.io/badge/-VMWare-0d0d0d?style=flat&logo=VMWare)![Nginx](https://img.shields.io/badge/-Nginx-0d0d0d?style=flat&logo=nginx)
```js
Universidad San Carlos de Guatemala 2022
Programador: Elder Anibal Pum Rojas
Carne: 201700761
Correo: ElderPum@gmail.com
```
---

### Universidad de San Carlos de Guatemala
### Primer Semestre 2024
### Descripción
Esta máquina virtual tiene la tarea de monitorear el estado de las máquinas virtuales y su desempeño a términos de recursos. Para esto, se utiliza un servicio desarrollado el cuál se encarga de recopilar y mostrar información importante, como la memoria RAM y el proceso del CPU, similar al Administrador de Tareas, pero para distribuciones Linux.

### Servicios Implementados
1. Docker: Utiliza imágenes docker que contienen la base de datos, el backend y el frontend, además de implementar docker compose para el servicio de contenedores.
2. C: Se utiliza el lenguaje C para el desarrollo de los módulos que se encargan de obtener desde el Kernel de Ubuntu los datos de la memoria RAM y el CPU.
3. Golang: Se utiliza una API desarrollada en Golang que se encarga de obtener los datos de los módulos, insertarlos en la base de datos y además, devolver los endpoints que el frontend consuma.
4. MySQL: Una base de datos persistente alojada en un contenedor docker que utiliza el lenguaje MySQL, con una única tabla donde se guarda la información.
5. ReactJS + Vite, Typescript, NodeJS: Se utiliza una aplicación web creada con la combinación de tecnologías de React y Vite para el rápido despliegue de la aplicación, Typescript como lenguaje de programación y NodeJS como servicio de paquetería de librerías y compilador general de la aplicación.
6. VMWare: Se utiliza como máquina virtual en la cuál, estará alojada nuestra aplicación con los distintos contenedores docker corriendo, justo acá es de donde se obtienen los datos de los procesos.
7. Nginx: Utilizado en la parte del frontend, sirve para poder conectarnos desde nuestra máquina principal al frontend alojado dentro de VMWare, así como también, permite conectarnos a los demás contenedores, sirve como puerta de enlace entre nuestra máquina local y la máquina virtual.
8. Ubuntu Server: Distro en formato ISO que se carga en VMWare y es el sistema operativo donde estará corriendo nuestra aplicación.

## Comandos Aplicados
### Módulos
Para la implementación de los módulos, es necesario instalar el compilador "GCC" y la herramienta "Make". Por lo cual se utilizan los siguiente comandos:

* GCC y Make

> sudo apt install build-essential

> sudo apt-get install manpages-dev

Para verificar si la instalación fue correcta, se usa lo siguiente:
> gcc --version

> make --version

Dentro de la carpeta '/module' se encuentran los módulos de CPU y RAM. Cada módulo se encuentra en su propia carpeta y dentro de cada carpeta se encuentra un archivo Makefile y un módulo C.
![](https://i.ibb.co/S04zJw3/image.png)
![](https://i.ibb.co/rw5QQXr/image.png)

Para poder ejecutar los módulos utilizamos lo siguiente:
* Ejecutar make:

> make all

* Limpiar archivos make:

> make clean

* Levantar make con sudo insmod:

> make run

* Detener make con sudo rmmod:

> make stop

### Backend
El backend se desarrolló en Golang, por lo que solamente se utilizó un único archivo llamado main.go, además, se sube a una imagen de docker utilizando su propio archivo dockerfile.

![](https://i.ibb.co/5Lh6pNw/image.png)

### Frontend
El frontend se realizó en ReactJS en conjunto con Vite para el despliegue rápido y eficiente de nuestra aplicación, igualmente, en conjunto con su respectivo archivo dockerfile para subir la imagen a Docker Hub.

Para levantar el frontend se utiliza lo siguiente:

* Para instalar el package.json

> npm i

* Para crear el archivo dist, que estará en nginx

> npm run build

* Para levantar el servidor de pruebas

> npm run dev

### Base de Datos
La base de datos se realizó en una imagen de Docker que contiene una base de datos en MySQL, donde únicamente existe una tabla que recopila información acerca de la RAM y del CPU.
![](https://i.ibb.co/6t7Jxz6/image.png)

### Imágenes y Contenedores de Docker
Para desplegar nuestra correcta aplicación, hicimos uso de imágenes publicadas dentro de Docker Hub ayudándonos de archivos dockerfile, además de utilizar un archivo docker compose que nos permite orquestar 3 contenedores de manera simultánea para nuestros procesos.
Los comandos a utilizar son:

* Para crear la imagen (independiente de qué servicio)

> docker build -f dockerfile -t epum/backend-so1-proyecto1:latest .

* Para pushear la imagen a Docker Hub (independiente de qué servicio)

> docker push epum/db-so1-proyecto1:latest

* Para levantar el archivo de docker-compose.yml

> docker-compose up