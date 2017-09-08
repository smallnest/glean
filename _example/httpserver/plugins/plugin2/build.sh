go build -ldflags "-pluginpath=plugin/hot-$(uuidgen)" -buildmode=plugin -o plugin2.so main.go
sha1sum plugin2.so > plugin2.sha1