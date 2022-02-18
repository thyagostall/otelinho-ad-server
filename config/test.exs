import Config

# Configure your database
#
# The MIX_TEST_PARTITION environment variable can be used
# to provide built-in test partitioning in CI environment.
# Run `mix help test` for more information.
config :otelinho_ad_server, OtelinhoAdServer.Repo,
  username: "otelinho",
  password: "devpassword",
  hostname: "localhost",
  database: "otelinho_ad_server_test#{System.get_env("MIX_TEST_PARTITION")}",
  pool: Ecto.Adapters.SQL.Sandbox,
  pool_size: 10

# We don't run a server during test. If one is required,
# you can enable the server option below.
config :otelinho_ad_server, OtelinhoAdServerWeb.Endpoint,
  http: [ip: {127, 0, 0, 1}, port: 4002],
  secret_key_base: "xNNIgAl3dDnXBwefMu2AhYZXsQBh1LrP3rYpHmSZ1r5Ipt8AmE5NW+VJFZHoI0c2",
  server: false

# In test we don't send emails.
config :otelinho_ad_server, OtelinhoAdServer.Mailer, adapter: Swoosh.Adapters.Test

# Print only warnings and errors during test
config :logger, level: :warn

# Initialize plugs at runtime for faster test compilation
config :phoenix, :plug_init_mode, :runtime
