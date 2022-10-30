using Server;
using TMPro;
using UnityEngine;

namespace UI
{
    public class LobbyItem : MonoBehaviour
    {
        [SerializeField] private TMP_Text lead;
        [SerializeField] private TMP_Text people;

        public void SetText(OnlinePlayer leadPlayer, uint curPeople, uint maxPeople)
        {
            lead.SetText("Lead: " + leadPlayer.ID);
            people.SetText(curPeople + " / " + maxPeople);
        }
    }
}