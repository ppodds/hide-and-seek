using System.Collections.Generic;
using Newtonsoft.Json;

namespace Server
{
    public class Lobby
    {
        public Lobby(uint id, OnlinePlayer lead, List<OnlinePlayer> players, byte currentPeople, byte maxPeople)
        {
            ID = id;
            Lead = lead;
            Players = players;
            CurrentPeople = currentPeople;
            MaxPeople = maxPeople;
        }

        public uint ID { get; set; }

        public OnlinePlayer Lead { get; set; }

        public List<OnlinePlayer> Players { get; set; }

        [JsonProperty("curPeople")]
        public byte CurrentPeople { get; set; }

        public byte MaxPeople { get; set; }
    }
}