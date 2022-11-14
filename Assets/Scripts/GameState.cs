using System;
using System.Collections.Generic;
using Protos;

[Serializable]
public class GameState
{
    public CharacterType Winner;
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