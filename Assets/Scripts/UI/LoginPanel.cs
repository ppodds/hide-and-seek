using TMPro;
using UnityEngine;

namespace UI
{
    public class LoginPanel : MonoBehaviour
    {
        [SerializeField] private TMP_InputField host;
        [SerializeField] private TMP_InputField tcpPort;
        [SerializeField] private TMP_InputField udpPort;

        public void Login()
        {
            var task = GameManager.Instance.ConnectToServer(host.text, int.Parse(tcpPort.text),
                int.Parse(udpPort.text));
            task.GetAwaiter().OnCompleted(() =>
            {
                if (!task.Result)
                {
                    Debug.Log("Login Failed, Please retry!");
                    return;
                }
                gameObject.SetActive(false);
                Debug.Log("Login Success");
            });
        }
    }
}