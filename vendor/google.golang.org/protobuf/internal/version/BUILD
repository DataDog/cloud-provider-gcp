load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "version",
    srcs = ["version.go"],
    importmap = "k8s.io/cloud-provider-gcp/vendor/google.golang.org/protobuf/internal/version",
    importpath = "google.golang.org/protobuf/internal/version",
    visibility = ["//vendor/google.golang.org/protobuf:__subpackages__"],
)
