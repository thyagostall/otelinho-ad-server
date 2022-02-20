defmodule OtelinhoAdServer.BeaconGenerator do
  def generate(event, campaign) do
    payload = encrypt(%{campaign_id: campaign.id})
    "http://localhost:4000/event/#{event}/#{payload}"
  end

  defp encrypt(val) do
    Plug.Crypto.encrypt(secret(), to_string(:beacon), val) |> Base.encode64()
  end

  def decrypt(ciphertext) do
    {:ok, decoded_ciphertext} = Base.decode64(ciphertext)
    {:ok, result} = Plug.Crypto.decrypt(secret(), to_string(:beacon), decoded_ciphertext)
    result
  end

  defp secret() do
    Application.get_env(:otelinho_ad_server, __MODULE__)[:secret]
  end
end
