load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "scheduling",
    srcs = [
        "interface.go",
        "priorityclass.go",
    ],
    importmap = "k8s.io/cloud-provider-gcp/vendor/k8s.io/client-go/informers/scheduling/v1",
    importpath = "k8s.io/client-go/informers/scheduling/v1",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/api/scheduling/v1:scheduling",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:meta",
        "//vendor/k8s.io/apimachinery/pkg/runtime",
        "//vendor/k8s.io/apimachinery/pkg/watch",
        "//vendor/k8s.io/client-go/informers/internalinterfaces",
        "//vendor/k8s.io/client-go/kubernetes",
        "//vendor/k8s.io/client-go/listers/scheduling/v1:scheduling",
        "//vendor/k8s.io/client-go/tools/cache",
    ],
)
