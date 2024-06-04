# gRPC API Example
## Cara Menjalankan
1. Clone repository ini
```bash
git clone git@github.com:fastcampus-backend-golang/grpc-api-example.git
```

2. Masuk ke direktori
```bash
cd grpc-api-example
```

3. Jalankan dengan perintah
```bash
go run .
```

## Konten
- main.go : file utama yang berisi konfigurasi server
- service.go : file yang berisi implementasi grpc api
- stock.proto : file protobuf rancangan api grpc
- proto: direktori yang berisi file protobuf yang sudah di-compile
- data: direktori yang berisi library untuk data saham palsu

## Cara Menggunakan
1. Jalankan server gRPC
2. Buka aplikasi desktop (Postman)[https://www.postman.com/]
3. Buat request baru, pilih gRPC
4. Pada bagian URL, masukkan url `grpc://localhost:50051`
5. Pada bagian Method, pilih `Use Server Reflection` jika reflection masih aktif atau `Import a .proto file`
6. Client sudah siap digunakan
 