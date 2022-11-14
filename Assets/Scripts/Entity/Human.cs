using IO.Net;
using SUPERCharacter;
using UnityEngine;

namespace Entity
{
    [RequireComponent(typeof(Animator), typeof(SUPERCharacterAIO), typeof(NetObject))]
    public class Human : MonoBehaviour
    {
        private static readonly int Velocity = Animator.StringToHash("velocity");
        private static readonly int Idle = Animator.StringToHash("idle");
        private static readonly int Running = Animator.StringToHash("running");
        private Animator _animator;
        private bool _cameraGenerated;
        private bool _isDead;
        private NetObject _netObject;
        private Rigidbody _rigidbody;
        private SUPERCharacterAIO _superCharacterAio;

        private void Awake()
        {
            _rigidbody = GetComponent<Rigidbody>();
            _animator = GetComponent<Animator>();
            _superCharacterAio = GetComponent<SUPERCharacterAIO>();
            _netObject = GetComponent<NetObject>();
            _isDead = false;
        }

        private void Update()
        {
            if (transform.position.y < -300)
            {
                GenerateSpectatorCamera();
                return;
            }

            if (_rigidbody.velocity.magnitude > 0.1)
            {
                _animator.SetFloat(Velocity, _rigidbody.velocity.magnitude);
                _animator.SetBool(Idle, false);
            }
            else
            {
                _animator.SetFloat(Velocity, 0);
                _animator.SetBool(Idle, true);
            }

            _animator.SetBool(Running, _superCharacterAio.isSprinting);
        }

        private void OnCollisionEnter(Collision collision)
        {
            if (_netObject.IsRemote)
                return;
            if (GenerateSpectatorCamera()) return;
            if (collision.gameObject.GetComponent<Ghost>() == null)
                return;
            var dir = (transform.position - collision.transform.position).normalized;
            _rigidbody.AddForce(new Vector3(dir.x, 0, dir.z) * 1500 + Vector3.up * 1000);
            _superCharacterAio.enableMovementControl = false;
            _isDead = true;
        }

        private bool GenerateSpectatorCamera()
        {
            if (_isDead && !_cameraGenerated)
            {
                _cameraGenerated = true;
                var t = transform;
                Instantiate(Resources.Load("Prefabs/Camera"), t.position,
                    t.rotation);
                _netObject.IsDead = true;
                return true;
            }

            return false;
        }
    }
}