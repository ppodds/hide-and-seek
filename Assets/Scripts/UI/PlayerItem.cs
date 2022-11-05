using Protos;
using TMPro;
using UnityEngine;

namespace UI
{
    public class PlayerItem : MonoBehaviour
    {
        [SerializeField] private TMP_Text lead;
        [SerializeField] private TMP_Text player;

        public void SetText(Lobby lobby, Player p)
        {
            lead.gameObject.SetActive(lobby.Lead.Id == p.Id);
            player.SetText(p.Id.ToString());
        }
    }
}