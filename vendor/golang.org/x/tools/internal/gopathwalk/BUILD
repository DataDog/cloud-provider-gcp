load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gopathwalk",
    srcs = ["walk.go"],
    importmap = "k8s.io/cloud-provider-gcp/vendor/golang.org/x/tools/internal/gopathwalk",
    importpath = "golang.org/x/tools/internal/gopathwalk",
    visibility = ["//vendor/golang.org/x/tools:__subpackages__"],
)
