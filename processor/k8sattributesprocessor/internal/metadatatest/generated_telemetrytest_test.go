// Code generated by mdatagen. DO NOT EDIT.

package metadatatest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/metadata"
)

func TestSetupTelemetry(t *testing.T) {
	testTel := SetupTelemetry()
	tb, err := metadata.NewTelemetryBuilder(
		testTel.NewTelemetrySettings(),
	)
	require.NoError(t, err)
	require.NotNil(t, tb)
	tb.OtelsvcK8sIPLookupMiss.Add(context.Background(), 1)
	tb.OtelsvcK8sNamespaceAdded.Add(context.Background(), 1)
	tb.OtelsvcK8sNamespaceDeleted.Add(context.Background(), 1)
	tb.OtelsvcK8sNamespaceUpdated.Add(context.Background(), 1)
	tb.OtelsvcK8sNodeAdded.Add(context.Background(), 1)
	tb.OtelsvcK8sNodeDeleted.Add(context.Background(), 1)
	tb.OtelsvcK8sNodeUpdated.Add(context.Background(), 1)
	tb.OtelsvcK8sPodAdded.Add(context.Background(), 1)
	tb.OtelsvcK8sPodDeleted.Add(context.Background(), 1)
	tb.OtelsvcK8sPodTableSize.Record(context.Background(), 1)
	tb.OtelsvcK8sPodUpdated.Add(context.Background(), 1)
	tb.OtelsvcK8sReplicasetAdded.Add(context.Background(), 1)
	tb.OtelsvcK8sReplicasetDeleted.Add(context.Background(), 1)
	tb.OtelsvcK8sReplicasetUpdated.Add(context.Background(), 1)

	testTel.AssertMetrics(t, []metricdata.Metrics{
		{
			Name:        "otelcol_otelsvc_k8s_ip_lookup_miss",
			Description: "Number of times pod by IP lookup failed.",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_namespace_added",
			Description: "Number of namespace add events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_namespace_deleted",
			Description: "Number of namespace delete events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_namespace_updated",
			Description: "Number of namespace update events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_node_added",
			Description: "Number of node add events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_node_deleted",
			Description: "Number of node delete events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_node_updated",
			Description: "Number of node update events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_pod_added",
			Description: "Number of pod add events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_pod_deleted",
			Description: "Number of pod delete events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_pod_table_size",
			Description: "Size of table containing pod info",
			Unit:        "1",
			Data: metricdata.Gauge[int64]{
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_pod_updated",
			Description: "Number of pod update events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_replicaset_added",
			Description: "Number of ReplicaSet add events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_replicaset_deleted",
			Description: "Number of ReplicaSet delete events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
		{
			Name:        "otelcol_otelsvc_k8s_replicaset_updated",
			Description: "Number of ReplicaSet update events received",
			Unit:        "1",
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{},
				},
			},
		},
	}, metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue())
	require.NoError(t, testTel.Shutdown(context.Background()))
}