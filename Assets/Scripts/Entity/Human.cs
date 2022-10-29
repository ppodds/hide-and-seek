using IO;
using SUPERCharacter;
using UnityEngine;

namespace Entity
{
    [RequireComponent(typeof(Animator), typeof(Rigidbody), typeof(SUPERCharacterAIO))]
    public class Human : MonoBehaviour
    {
        private Rigidbody _rigidbody;
        private Animator _animator;
        private SUPERCharacterAIO _superCharacterAio;
        private static readonly int Velocity = Animator.StringToHash("velocity");
        private static readonly int Idle = Animator.StringToHash("idle");
        private static readonly int Running = Animator.StringToHash("running");

        private void Awake()
        {
            _rigidbody = GetComponent<Rigidbody>();
            _animator = GetComponent<Animator>();
            _superCharacterAio = GetComponent<SUPERCharacterAIO>();
        }

        private void Update()
        {
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
    }
}