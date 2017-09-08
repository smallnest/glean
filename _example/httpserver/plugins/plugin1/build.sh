go build -ldflags "-pluginpath=plugin/hot-$(uuidgen)" -buildmode=plugin -o plugin1.so main.go
sha1sum plugin1.so > plugin1.sha1