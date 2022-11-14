using TMPro;
using UnityEngine;

namespace UI.Toast
{
    [RequireComponent(typeof(TMP_Text))]
    public class ToastItem : MonoBehaviour
    {
        private TMP_Text _text;

        private void Awake()
        {
            _text = GetComponent<TMP_Text>();
        }

        public void ShowText(string text)
        {
            _text.SetText(text);
        }
    }
}