# .PHONY tells make that these rules do not correspond to actual files
.PHONY: protobuf

codegen:
  oapi-codegen -package qbittorrent  openapi.yaml > ./qbittorrent/qbittorrent.gen.go

