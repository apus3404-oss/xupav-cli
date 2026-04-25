# mycli TUI Nasıl Çalıştırılır

## Adım 1: Terminal Aç
Windows PowerShell veya Git Bash'i aç

## Adım 2: Proje Dizinine Git
```bash
cd D:\my-cli
```

## Adım 3: Chat Komutunu Çalıştır
```bash
.\bin\mycli.exe chat
```

## Ne Göreceksin:
- 🤖 Başlık: "mycli | Messages: 1"
- Hoş geldin mesajı (System mesajı)
- Alt kısımda input kutusu
- En altta yardım metni: "Enter: send | Ctrl+C: quit"

## Nasıl Kullanılır:
1. Input kutusuna mesaj yaz
2. Enter'a bas (mesaj gönderilir)
3. Ctrl+C ile çık

## Not:
Şu anda Python bridge başlatılıyor ama gerçek API key olmadığı için 
AI yanıt vermeyecek. Ama TUI'nin görünümünü ve input handling'i görebilirsin!

## Alternatif: Sadece Görünümü Görmek İçin
Eğer sadece nasıl göründüğünü görmek istiyorsan:
```bash
.\bin\mycli.exe config show
```
Bu komut TUI olmadan config'i gösterir.
