load("@rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "go_links_server",
    srcs = [
        "server.go"
    ],
    cgo = True,
    deps = [
        "//middleware",
        "//handlers",
        "@com_github_adrg_xdg//:xdg",
        "@com_github_mattn_go_sqlite3//:go-sqlite3"
    ],
)
