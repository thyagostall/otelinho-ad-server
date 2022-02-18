defmodule OtelinhoAdServerWeb.Beacons do
  use OtelinhoAdServerWeb, :controller

  def event(conn, _params) do
    gif_data = <<71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0, 255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 1, 68, 0, 59>>
    send_resp(conn, 200, gif_data)
  end
end
