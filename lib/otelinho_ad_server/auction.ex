defmodule OtelinhoAdServer.Auction do

  def run_auction(campaigns) do
    campaigns
    |> sort_by_max_bid()
    |> choose_campaign()
  end

  defp sort_by_max_bid(campaigns) do
    campaigns
    |> Enum.sort_by(fn c -> c.max_bid end, {:desc, Decimal})
  end

  defp choose_campaign([]), do: nil
  defp choose_campaign([campaign]), do: %{campaign | max_bid: Decimal.new("0.25")}
  defp choose_campaign([first, second | _]) do
    %{first | max_bid: Decimal.add(second.max_bid, Decimal.new("0.01"))}
  end
end
