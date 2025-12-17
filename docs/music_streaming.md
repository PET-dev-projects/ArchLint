# Music Streaming Reference Architecture

This document visualizes `examples/music_streaming.yaml` with a Mermaid diagram so you can review the end-to-end topology at a glance.

## Mermaid diagram

```mermaid
flowchart LR
  subgraph Playback[Playback & Streaming]
    playback_api[playback-api]
    session_service[session-service]
    session_store[(session-store)]
    streaming_control[streaming-control]
    cdn_publisher[cdn-publisher]
  end

  subgraph Catalog[Catalog & Library]
    catalog_api[catalog-api]
    catalog_workflow[catalog-workflow]
    catalog_repo[catalog-repo]
    catalog_db[(catalog-db)]
    search_indexer[search-indexer]
    licensing_proxy[licensing-proxy]
  end

  subgraph Personalization[Personalization & ML]
    listening_ingest[listening-ingest]
    listening_store[(listening-store)]
    feature_service[feature-service]
    feature_store[(feature-store)]
    model_service[model-service]
    recommendation_api[recommendation-api]
    experiment_service[experiment-service]
  end

  subgraph Monetization[Monetization & Ads]
    billing_api[billing-api]
    billing_ledger[billing-ledger]
    billing_db[(billing-db)]
    billing_adapter[billing-adapter]
    ads_orchestrator[ads-orchestrator]
    ads_repo[ads-repo]
    ads_db[(ads-db)]
    reporting_service[reporting-service]
    ad_broker[ad-broker]
  end

  subgraph Platform[Platform Services]
    identity_service[identity-service]
    entitlement_service[entitlement-service]
    entitlement_db[(entitlement-db)]
    audit_log[audit-log]
    audit_db[(audit-db)]
    notification_relay[notification-relay]
    monitoring_pipeline[monitoring-pipeline]
    data_warehouse[(data-warehouse)]
  end

  global_cdn{{global-cdn}}
  rights_societies{{rights-societies}}
  payment_gateway{{payment-gateway}}
  ad_network{{ad-network}}
  email_provider{{email-provider}}

  playback_api -->|sync| session_service
  session_service -->|db| session_store
  session_service -->|sync| streaming_control
  streaming_control -->|async| cdn_publisher
  cdn_publisher -->|async| global_cdn
  playback_api -->|sync| entitlement_service
  streaming_control -->|async| listening_ingest

  catalog_api -->|sync| catalog_workflow
  catalog_workflow -->|sync| catalog_repo
  catalog_repo -->|db| catalog_db
  search_indexer -->|sync| catalog_repo
  catalog_api -->|async| search_indexer
  licensing_proxy -->|sync| catalog_workflow
  licensing_proxy -->|async| rights_societies

  listening_ingest -->|db| listening_store
  listening_ingest -->|async| feature_service
  feature_service -->|db| feature_store
  model_service -->|sync| feature_service
  recommendation_api -->|sync| model_service
  experiment_service -->|sync| recommendation_api
  recommendation_api -->|async| notification_relay

  billing_api -->|sync| billing_ledger
  billing_ledger -->|db| billing_db
  billing_ledger -->|sync| billing_adapter
  billing_adapter -->|sync| payment_gateway
  ads_orchestrator -->|sync| ads_repo
  ads_repo -->|db| ads_db
  reporting_service -->|db| ads_db
  ads_orchestrator -->|sync| ad_broker
  ad_broker -->|async| ad_network
  billing_ledger -->|async| audit_log

  identity_service -->|sync| entitlement_service
  entitlement_service -->|db| entitlement_db
  audit_log -->|db| audit_db
  monitoring_pipeline -->|db| data_warehouse
  audit_log -->|async| monitoring_pipeline
  notification_relay -->|async| email_provider
```

> Tip: render the diagram locally with a Markdown viewer that supports Mermaid, or copy the snippet to [mermaid.live](https://mermaid.live) for quick experimentation.
