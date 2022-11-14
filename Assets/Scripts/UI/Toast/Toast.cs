using System;
using System.Threading.Tasks;
using UnityEngine;

namespace UI.Toast
{
    public class Toast : MonoBehaviour
    {
        [SerializeField] private GameObject toastItem;

        public void PushToast(string message)
        {
            var item = Instantiate(toastItem, transform).GetComponent<ToastItem>();
            item.ShowText(message);
            Task.Delay(new TimeSpan(0, 0, 3)).GetAwaiter().OnCompleted(() => { Destroy(item.gameObject); });
        }

        public void PushToast(string message, TimeSpan duration)
        {
            var item = Instantiate(toastItem, transform).GetComponent<ToastItem>();
            item.ShowText(message);
            Task.Delay(duration).GetAwaiter().OnCompleted(() => { Destroy(item.gameObject); });
        }
    }
}