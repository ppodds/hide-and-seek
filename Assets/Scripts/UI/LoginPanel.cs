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
            GameManager.Instance.Server = new Server
            {
                Host = host.text,
                TcpPort = int.Parse(tcpPort.text),
                UdpPort = int.Parse(udpPort.text)
            };
            var task = GameManager.Instance.ConnectToServer();
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