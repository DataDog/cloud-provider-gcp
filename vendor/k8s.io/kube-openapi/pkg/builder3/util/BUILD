load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "util",
    srcs = ["util.go"],
    importmap = "k8s.io/cloud-provider-gcp/vendor/k8s.io/kube-openapi/pkg/builder3/util",
    importpath = "k8s.io/kube-openapi/pkg/builder3/util",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/kube-openapi/pkg/schemamutation",
        "//vendor/k8s.io/kube-openapi/pkg/validation/spec",
    ],
)
