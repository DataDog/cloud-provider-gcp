load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "recognizer",
    srcs = ["recognizer.go"],
    importmap = "k8s.io/cloud-provider-gcp/vendor/k8s.io/apimachinery/pkg/runtime/serializer/recognizer",
    importpath = "k8s.io/apimachinery/pkg/runtime/serializer/recognizer",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/runtime",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema",
    ],
)
