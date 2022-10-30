using UnityEngine;

namespace CameraController
{
    public class WelcomeCamera : MonoBehaviour
    {
        [SerializeField] private float rotateDegreePerSecond;
        private Camera _camera;

        private void Awake()
        {
            _camera = GetComponent<Camera>();
        }

        private void Update()
        {
            _camera.transform.Rotate(0, rotateDegreePerSecond * Time.deltaTime, 0);
        }
    }
}