load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "errors",
    srcs = [
        "doc.go",
        "errors.go",
    ],
    importmap = "k8s.io/cloud-provider-gcp/vendor/k8s.io/apimachinery/pkg/api/errors",
    importpath = "k8s.io/apimachinery/pkg/api/errors",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:meta",
        "//vendor/k8s.io/apimachinery/pkg/runtime",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema",
        "//vendor/k8s.io/apimachinery/pkg/util/validation/field",
    ],
)
