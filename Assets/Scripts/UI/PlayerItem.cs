using Server;
using TMPro;
using UnityEngine;

namespace UI
{
    public class PlayerItem : MonoBehaviour
    {
        [SerializeField] private TMP_Text lead;
        [SerializeField] private TMP_Text player;

        public void SetText(Lobby lobby, OnlinePlayer p)
        {
            lead.gameObject.SetActive(lobby.Lead.ID == p.ID);
            player.SetText(p.ID.ToString());
        }
    }
}