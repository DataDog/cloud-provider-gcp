load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go-strcase",
    srcs = [
        "camel.go",
        "doc.go",
        "helper.go",
        "kebab.go",
        "snake.go",
    ],
    importmap = "k8s.io/cloud-provider-gcp/vendor/github.com/stoewer/go-strcase",
    importpath = "github.com/stoewer/go-strcase",
    visibility = ["//visibility:public"],
)
