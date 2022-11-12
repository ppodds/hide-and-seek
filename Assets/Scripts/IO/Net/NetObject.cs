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
                var t = transform;
                var position = t.position;
                var rotation = t.rotation.eulerAngles;
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
                    Rotation = new Vector3
                    {
                        X = rotation.x,
                        Y = rotation.y,
                        Z = rotation.z
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

        public void SetVelocity(Vector3 v)
        {
            _rigidbody.velocity = new UnityEngine.Vector3(v.X, v.Y, v.Z);
        }

        public void SetRotation(Vector3 v)
        {
            transform.rotation = Quaternion.Euler(v.X, v.Y, v.Z);
        }

        public void SetPosition(Vector3 v)
        {
            transform.position = new UnityEngine.Vector3(v.X, v.Y, v.Z);
        }
    }
}