defmodule OtelinhoAdServer.Repo do
  use Ecto.Repo,
    otp_app: :otelinho_ad_server,
    adapter: Ecto.Adapters.Postgres
end
