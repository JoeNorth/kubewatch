package config

var yannotated = `# Handlers know how to send notifications to specific services.
handler:
  slack:
    # Slack "legacy" API token.
    token: ""
    # Slack channel.
    channel: ""
    # Title of the message.
    title: ""
  slackwebhook:
    # Slack channel.
    channel: ""
    # Slack Username.
    username: ""
    # Slack Emoji.
    emoji: ""
    # Slack Webhook Url.
    slackwebhookurl: ""
  hipchat:
    # Hipchat token.
    token: ""
    # Room name.
    room: ""
    # URL of the hipchat server.
    url: ""
  mattermost:
    room: ""
    url: ""
    username: ""
  flock:
    # URL of the flock API.
    url: ""
  webhook:
    # Webhook URL.
    url: ""
    cert: ""
    tlsskip: false
  cloudevent:
    url: ""
  msteams:
    # MSTeams API Webhook URL.
    webhookurl: ""
  smtp:
    # Destination e-mail address.
    to: ""
    # Sender e-mail address .
    from: ""
    # Smarthost, aka "SMTP server"; address of server used to send email.
    smarthost: ""
    # Subject of the outgoing emails.
    subject: ""
    # Extra e-mail headers to be added to all outgoing messages.
    headers: {}
    # Authentication parameters.
    auth:
      # Username for PLAN and LOGIN auth mechanisms.
      username: ""
      # Password for PLAIN and LOGIN auth mechanisms.
      password: ""
      # Identity for PLAIN auth mechanism
      identity: ""
      # Secret for CRAM-MD5 auth mechanism
      secret: ""
    # If "true" forces secure SMTP protocol (AKA StartTLS).
    requireTLS: false
    # SMTP hello field (optional)
    hello: ""
  lark:
    # Webhook URL.
    webhookurl: ""
  eventbridge:
    # EventBridge EndpointId (optional)
    endpointId: ""
    # EKS Cluster ARN. Used for EventBridge Event resource (optional)
    clusterArn: ""
    # EventBridge EventBusName (optional)
    eventBusName: ""
# Resources to watch.
resource:
  deployment: false
  rc: false
  rs: false
  ds: false
  statefulset: false
  svc: false
  po: false
  job: false
  node: false
  clusterrole: false
  clusterrolebinding: false
  sa: false
  pv: false
  ns: false
  secret: false
  configmap: false
  ing: false
  hpa: false
  event: false
  coreevent: false
# For watching specific namespace, leave it empty for watching all.
# this config is ignored when watching namespaces
namespace: ""
`
