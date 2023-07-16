<p align="center">
<img src="https://cdn.jarvans.com/blog/2023/202306042109858.jpg" width="200px" style="border-radius: 50%;"/>
<br>
<p align="center">
 <img src="https://img.shields.io/github/stars/jarvanstack/pn" />
 <img src="https://img.shields.io/github/issues/jarvanstack/pn" />
 <img src="https://img.shields.io/github/forks/jarvanstack/pn" />
</p>
</p>

[简体中文](./README.md) 

## PasswordNote

[link](https://web.jarvans.com/pn_client)

PasswordNote is a password management tool that generates various types of passwords with a single click. It employs client-side encryption, where the server only stores the hash value of the master password. This ensures that only you know your password, and only you can access your data. PasswordNote supports one-click deployment of your own server for data privatization.

Features:

1. One-click generation of various types of passwords.
2. Supports online and offline storage modes.
3. The server only stores the hash value of the master password, ensuring that only you know your password.
4. The server only stores the client-side encrypted notes, ensuring that only you can access your data.
5. Supports one-click deployment of your own server for data privatization.

## Screenshots

![Generate Password](https://cdn.jarvans.com/blog/2023/202306041848001.png)

![Home Page](https://cdn.jarvans.com/blog/2023/202306041852072.png)

## Method 1: Development Mode

```bash
cd pn_server
# Compile and run
make run

```

* Requires Golang environment.

## Method 2: Docker Mode

```bash
cd pn_server
# build
make docker.build
# run
make up
```

* Requires Docker environment.

## Method 3: Docker Image

TODO

## Method 4: One-Click Server Deployment

TODO
