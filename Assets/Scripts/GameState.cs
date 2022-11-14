using System;
using System.Collections.Generic;

[Serializable]
public class GameState
{
    public readonly uint Id;
    public readonly Dictionary<uint, PlayerState> Players;

    public GameState(uint id)
    {
        Id = id;
        Players = new Dictionary<uint, PlayerState>();
        InGame = true;
    }

    public bool InGame { get; set; }
}