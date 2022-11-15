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
        [field: SerializeField] public uint PlayerId { get; set; }
        [field: SerializeField] public bool IsDead { get; set; }

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
                var task = _udpClient.UpdatePlayer(new Character
                {
                    Dead = IsDead,
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
                task.GetAwaiter().OnCompleted(() =>
                {
                    if (IsDead)
                        Destroy(gameObject);
                });
            }
            else
            {
                if (GameManager.Instance.GameState == null)
                {
                    Destroy(gameObject);
                    return;
                }

                var character = GameManager.Instance.GameState.Players[PlayerId].Player.Character;
                IsDead = character.Dead;
                if (IsDead)
                {
                    Destroy(gameObject);
                    return;
                }

                SetPosition(character.Pos);
                SetRotation(character.Rotation);
                SetVelocity(character.Velocity);
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