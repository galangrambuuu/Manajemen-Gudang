package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Barang struct {
	KodeBarang string
	NamaBarang string
	Jumlah     int
	Harga      float64
}

type Transaksi struct {
	KodeBarang string
	NamaBarang string
	Jumlah     int
	Harga      float64
	Tipe       string
}

type Node struct {
	Data Transaksi
	Next *Node
}

type LinkedList struct {
	Head *Node
}

func (ll *LinkedList) Add(data Transaksi) {
	newNode := &Node{Data: data}
	if ll.Head == nil {
		ll.Head = newNode
	} else {
		current := ll.Head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newNode
	}
}

func (ll *LinkedList) ProsesAntrean() {
	if ll.Head == nil {
		fmt.Println("Antrean kosong, tidak ada yang bisa diproses.")
		return
	}

	transaksi := ll.Head.Data

	barang := cariBarang(transaksi.KodeBarang)
	if barang != nil {

		barang.Jumlah += transaksi.Jumlah
		fmt.Printf("Barang %s ditambahkan ke stok. Jumlah stok sekarang: %d\n", barang.NamaBarang, barang.Jumlah)
	} else {

		barangBaru := Barang{
			KodeBarang: transaksi.KodeBarang,
			NamaBarang: transaksi.NamaBarang,
			Jumlah:     transaksi.Jumlah,
			Harga:      transaksi.Harga,
		}
		daftarBarang = append(daftarBarang, barangBaru)
		fmt.Printf("Barang %s ditambahkan ke stok baru. Jumlah stok: %d\n", barangBaru.NamaBarang, barangBaru.Jumlah)
	}

	ll.Head = ll.Head.Next
	fmt.Println("Antrean pertama diproses dan dihapus.")
}

func (ll *LinkedList) Display() {
	if ll.Head == nil {
		fmt.Println("Antrean kosong.")
		return
	}
	fmt.Println("Antrean Pembelian:")
	current := ll.Head
	for current != nil {
		fmt.Printf("Kode: %s, Nama: %s, Jumlah: %d, Harga: %.2f\n",
			current.Data.KodeBarang, current.Data.NamaBarang, current.Data.Jumlah, current.Data.Harga)
		current = current.Next
	}
}

var daftarBarang []Barang
var stackPenjualan []Transaksi
var queuePembelianLL LinkedList
var fileName = "data_barang.json"
var filePenjualan = "riwayat_penjualan.json"

func tambahBarang(kodeBarang, namaBarang string, jumlah int, harga float64) {
	barang := Barang{
		KodeBarang: kodeBarang,
		NamaBarang: namaBarang,
		Jumlah:     jumlah,
		Harga:      harga,
	}
	daftarBarang = append(daftarBarang, barang)
	fmt.Println("Barang berhasil ditambahkan.")
	simpanDataKeFile()
}

func jualBarang(kodeBarang string, jumlah int) {
	for i, barang := range daftarBarang {
		if barang.KodeBarang == kodeBarang {
			if barang.Jumlah >= jumlah {
				daftarBarang[i].Jumlah -= jumlah
				transaksi := Transaksi{kodeBarang, barang.NamaBarang, jumlah, barang.Harga, "penjualan"}
				stackPenjualan = append(stackPenjualan, transaksi)
				fmt.Println("Barang berhasil dijual.")
				simpanDataKeFile()
				simpanRiwayatPenjualan()
				return
			} else {
				fmt.Println("Stok tidak mencukupi.")
				return
			}
		}
	}
	fmt.Println("Barang tidak ditemukan.")
}

func simpanDataKeFile() {
	data, err := json.Marshal(daftarBarang)
	if err != nil {
		fmt.Println("Gagal menyimpan data barang:", err)
		return
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Gagal menyimpan data ke file:", err)
	}
}

func simpanRiwayatPenjualan() {
	data, err := json.Marshal(stackPenjualan)
	if err != nil {
		fmt.Println("Gagal menyimpan riwayat penjualan:", err)
		return
	}
	err = os.WriteFile(filePenjualan, data, 0644)
	if err != nil {
		fmt.Println("Gagal menyimpan riwayat penjualan ke file:", err)
	}
}

func sortBarang(by string) {
	switch by {
	case "nama_asc":
		sort.Slice(daftarBarang, func(i, j int) bool {
			return daftarBarang[i].NamaBarang < daftarBarang[j].NamaBarang
		})
	case "harga_asc":
		sort.Slice(daftarBarang, func(i, j int) bool {
			return daftarBarang[i].Harga < daftarBarang[j].Harga
		})
	case "kode_asc":
		sort.Slice(daftarBarang, func(i, j int) bool {
			return daftarBarang[i].KodeBarang < daftarBarang[j].KodeBarang
		})
	}
	fmt.Println("Barang berhasil diurutkan.")
	displayBarang()
}

func displayBarang() {
	if len(daftarBarang) == 0 {
		fmt.Println("Daftar barang kosong.")
		return
	}
	fmt.Println("\nDaftar Barang yang Dimiliki:")
	for _, barang := range daftarBarang {
		fmt.Printf("Kode: %s, Nama: %s, Jumlah: %d, Harga: %.2f\n",
			barang.KodeBarang, barang.NamaBarang, barang.Jumlah, barang.Harga)
	}
}

func cariBarang(kodeBarang string) *Barang {
	for _, barang := range daftarBarang {
		if barang.KodeBarang == kodeBarang {
			return &barang
		}
	}
	return nil
}

func rekapPenjualan() {
	if len(stackPenjualan) == 0 {
		fmt.Println("Tidak ada riwayat penjualan.")
		return
	}
	fmt.Println("\nRekap Penjualan:")
	for i := len(stackPenjualan) - 1; i >= 0; i-- {
		transaksi := stackPenjualan[i]
		fmt.Printf("Kode: %s, Nama: %s, Jumlah: %d, Harga: %.2f, Tipe: %s\n",
			transaksi.KodeBarang, transaksi.NamaBarang, transaksi.Jumlah, transaksi.Harga, transaksi.Tipe)
	}
}

func main() {
	muatDataDariFile()
	var pilihan int
	for {
		fmt.Println("\n=== Manajemen Gudang ===")
		fmt.Println("1. Tambah Barang")
		fmt.Println("2. Jual Barang")
		fmt.Println("3. Urutkan Barang")
		fmt.Println("4. Cari Barang")
		fmt.Println("5. Kelola Antrean Pembelian ")
		fmt.Println("6. Tampilkan Barang yang Dimiliki")
		fmt.Println("7. Rekap Penjualan")
		fmt.Println("8. Keluar")
		fmt.Print("Pilih menu: ")
		fmt.Scanln(&pilihan)

		switch pilihan {
		case 1:
			kodeBarang := inputString("Masukkan Kode Barang: ")
			namaBarang := inputString("Masukkan Nama Barang: ")
			jumlah := inputJumlah()
			harga := inputHarga()
			tambahBarang(kodeBarang, namaBarang, jumlah, harga)
		case 2:
			kodeBarang := inputString("Masukkan Kode Barang yang Akan Dijual: ")
			jumlah := inputJumlah()
			jualBarang(kodeBarang, jumlah)
		case 3:
			fmt.Println("Urutkan Barang:")
			fmt.Println("1. Berdasarkan Nama ")
			fmt.Println("2. Berdasarkan Harga ")
			fmt.Println("3. Berdasarkan Kode")
			var subPilihan int
			fmt.Scanln(&subPilihan)
			switch subPilihan {
			case 1:
				sortBarang("nama_asc")
			case 2:
				sortBarang("harga_asc")
			case 3:
				sortBarang("kode_asc")
			default:
				fmt.Println("Pilihan tidak valid.")
			}
		case 4:
			kodeBarang := inputString("Masukkan Kode Barang: ")
			barang := cariBarang(kodeBarang)
			if barang != nil {
				fmt.Printf("Ditemukan: %s, Jumlah: %d, Harga: %.2f\n", barang.NamaBarang, barang.Jumlah, barang.Harga)
			} else {
				fmt.Println("Barang tidak ditemukan.")
			}
		case 5:
			fmt.Println("Kelola Antrean Pembelian:")
			fmt.Println("1. Tambah Antrean")
			fmt.Println("2. Proses Antrean Pembelian")
			fmt.Println("3. Tampilkan Antrean")
			var subPilihan int
			fmt.Scanln(&subPilihan)
			switch subPilihan {
			case 1:
				kodeBarang := inputString("Masukkan Kode Barang untuk Dibeli: ")
				namaBarang := inputString("Masukkan Nama Barang: ")
				jumlah := inputJumlah()
				harga := inputHarga()
				queuePembelianLL.Add(Transaksi{kodeBarang, namaBarang, jumlah, harga, "pembelian"})
				fmt.Println("Antrean pembelian ditambahkan.")
			case 2:
				queuePembelianLL.ProsesAntrean()
			case 3:
				queuePembelianLL.Display()
			default:
				fmt.Println("Pilihan tidak valid.")
			}
		case 6:
			displayBarang()
		case 7:
			rekapPenjualan()
		case 8:
			fmt.Println("Terima kasih!")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

func inputString(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

func inputJumlah() int {
	var jumlah int
	fmt.Print("Masukkan Jumlah: ")
	fmt.Scanln(&jumlah)
	return jumlah
}

func inputHarga() float64 {
	var harga float64
	fmt.Print("Masukkan Harga: ")
	fmt.Scanln(&harga)
	return harga
}
func muatDataDariFile() {
	data, err := os.ReadFile(fileName)
	if err == nil {
		json.Unmarshal(data, &daftarBarang)
	}
}
