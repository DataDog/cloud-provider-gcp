load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "apps",
    srcs = ["interface.go"],
    importmap = "k8s.io/cloud-provider-gcp/vendor/k8s.io/client-go/informers/apps",
    importpath = "k8s.io/client-go/informers/apps",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/client-go/informers/apps/v1:apps",
        "//vendor/k8s.io/client-go/informers/apps/v1beta1",
        "//vendor/k8s.io/client-go/informers/apps/v1beta2",
        "//vendor/k8s.io/client-go/informers/internalinterfaces",
    ],
)
