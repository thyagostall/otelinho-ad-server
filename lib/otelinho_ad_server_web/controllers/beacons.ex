defmodule OtelinhoAdServerWeb.Beacons do
  use OtelinhoAdServerWeb, :controller

  alias OtelinhoAdServer.Beacon
  alias OtelinhoAdServer.Repo
  alias OtelinhoAdServer.BeaconGenerator

  def event(conn, %{"event_type" => event_type, "event_data" => event_data}) do
    Task.async(fn -> process_beacon(event_type, event_data) end)

    gif_data = <<71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0, 255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 1, 68, 0, 59>>
    send_resp(conn, 200, gif_data)
  end

  defp process_beacon(event_type, event_data) do
    beacon = BeaconGenerator.decrypt(event_data)

    %Beacon{}
    |> Beacon.changeset(%{campaign_id: beacon.campaign_id, event: event_type, timestamp: DateTime.now!("Etc/UTC")})
    |> Repo.insert!()
  end
end
