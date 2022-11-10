using Protos;
using UnityEngine;
using Vector3 = Protos.Vector3;

namespace IO.Net
{
    [RequireComponent(typeof(Rigidbody))]
    public class NetObject : MonoBehaviour
    {
        private Rigidbody _rigidbody;
        private GameUdpClient _udpClient;
        [field: SerializeField] public bool IsRemote { get; set; }

        private void Awake()
        {
            _udpClient = new GameUdpClient(GameManager.Instance.Server.Host, GameManager.Instance.Server.UdpPort);
            _udpClient.Connect();
            _rigidbody = GetComponent<Rigidbody>();
        }

        private void Update()
        {
            if (!IsRemote)
            {
                var position = _rigidbody.position;
                var velocity = _rigidbody.velocity;
                _udpClient.UpdatePlayer(new Character
                {
                    Dead = false,
                    Velocity = new Vector3
                    {
                        X = velocity.x,
                        Y = velocity.y,
                        Z = velocity.z
                    },
                    Pos = new Vector3
                    {
                        X = position.x,
                        Y = position.y,
                        Z = position.z
                    }
                });
            }
        }
    }
}