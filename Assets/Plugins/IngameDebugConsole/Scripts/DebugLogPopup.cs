using System.Collections;
using UnityEngine;
using UnityEngine.EventSystems;
using UnityEngine.UI;
#if UNITY_EDITOR && UNITY_2021_1_OR_NEWER
using Screen = UnityEngine.Device.Screen; // To support Device Simulator on Unity 2021.1+
#endif

// Manager class for the debug popup
namespace IngameDebugConsole
{
    public class DebugLogPopup : MonoBehaviour, IPointerClickHandler, IBeginDragHandler, IDragHandler, IEndDragHandler
    {
        // Background image that will change color to indicate an alert
        private Image backgroundImage;

        // Canvas group to modify visibility of the popup
        private CanvasGroup canvasGroup;

        // Dimensions of the popup divided by 2
        private Vector2 halfSize;

        private bool isPopupBeingDragged;

        // Coroutines for simple code-based animations
        private IEnumerator moveToPosCoroutine;

        // Number of new debug entries since the log window has been closed
        private int newInfoCount, newWarningCount, newErrorCount;

        private Color normalColor;
        private Vector2 normalizedPosition;
        private RectTransform popupTransform;

        private void Awake()
        {
            popupTransform = (RectTransform)transform;
            backgroundImage = GetComponent<Image>();
            canvasGroup = GetComponent<CanvasGroup>();

            normalColor = backgroundImage.color;

            halfSize = popupTransform.sizeDelta * 0.5f;

            var pos = popupTransform.anchoredPosition;
            if (pos.x != 0f || pos.y != 0f)
                normalizedPosition = pos.normalized; // Respect the initial popup position set in the prefab
            else
                normalizedPosition = new Vector2(0.5f, 0f); // Right edge by default
        }

        private void Reset()
        {
            newInfoCount = 0;
            newWarningCount = 0;
            newErrorCount = 0;

            newInfoCountText.text = "0";
            newWarningCountText.text = "0";
            newErrorCountText.text = "0";

            backgroundImage.color = normalColor;
        }

        public void OnBeginDrag(PointerEventData data)
        {
            isPopupBeingDragged = true;

            // If a smooth movement animation is in progress, cancel it
            if (moveToPosCoroutine != null)
            {
                StopCoroutine(moveToPosCoroutine);
                moveToPosCoroutine = null;
            }
        }

        // Reposition the popup
        public void OnDrag(PointerEventData data)
        {
            Vector2 localPoint;
            if (RectTransformUtility.ScreenPointToLocalPointInRectangle(debugManager.canvasTR, data.position,
                    data.pressEventCamera, out localPoint))
                popupTransform.anchoredPosition = localPoint;
        }

        // Smoothly translate the popup to the nearest edge
        public void OnEndDrag(PointerEventData data)
        {
            isPopupBeingDragged = false;
            UpdatePosition(false);
        }

        // Popup is clicked
        public void OnPointerClick(PointerEventData data)
        {
            // Hide the popup and show the log window
            if (!isPopupBeingDragged)
                debugManager.ShowLogWindow();
        }

        public void NewLogsArrived(int newInfo, int newWarning, int newError)
        {
            if (newInfo > 0)
            {
                newInfoCount += newInfo;
                newInfoCountText.text = newInfoCount.ToString();
            }

            if (newWarning > 0)
            {
                newWarningCount += newWarning;
                newWarningCountText.text = newWarningCount.ToString();
            }

            if (newError > 0)
            {
                newErrorCount += newError;
                newErrorCountText.text = newErrorCount.ToString();
            }

            if (newErrorCount > 0)
                backgroundImage.color = alertColorError;
            else if (newWarningCount > 0)
                backgroundImage.color = alertColorWarning;
            else
                backgroundImage.color = alertColorInfo;
        }

        // A simple smooth movement animation
        private IEnumerator MoveToPosAnimation(Vector2 targetPos)
        {
            var modifier = 0f;
            var initialPos = popupTransform.anchoredPosition;

            while (modifier < 1f)
            {
                modifier += 4f * Time.unscaledDeltaTime;
                popupTransform.anchoredPosition = Vector2.Lerp(initialPos, targetPos, modifier);

                yield return null;
            }
        }

        // Hides the log window and shows the popup
        public void Show()
        {
            canvasGroup.blocksRaycasts = true;
            canvasGroup.alpha = 1f;

            // Reset the counters
            Reset();

            // Update position in case resolution was changed while the popup was hidden
            UpdatePosition(true);
        }

        // Hide the popup
        public void Hide()
        {
            canvasGroup.blocksRaycasts = false;
            canvasGroup.alpha = 0f;

            isPopupBeingDragged = false;
        }

        // There are 2 different spaces used in these calculations:
        // RectTransform space: raw anchoredPosition of the popup that's in range [-canvasSize/2, canvasSize/2]
        // Safe area space: Screen.safeArea space that's in range [safeAreaBottomLeft, safeAreaTopRight] where these corner positions
        //                  are all positive (calculated from bottom left corner of the screen instead of the center of the screen)
        public void UpdatePosition(bool immediately)
        {
            var canvasRawSize = debugManager.canvasTR.rect.size;

            // Calculate safe area bounds
            var canvasWidth = canvasRawSize.x;
            var canvasHeight = canvasRawSize.y;

            var canvasBottomLeftX = 0f;
            var canvasBottomLeftY = 0f;

            if (debugManager.popupAvoidsScreenCutout)
            {
#if UNITY_2017_2_OR_NEWER && ( UNITY_EDITOR || UNITY_ANDROID || UNITY_IOS )
                var safeArea = Screen.safeArea;

                var screenWidth = Screen.width;
                var screenHeight = Screen.height;

                canvasWidth *= safeArea.width / screenWidth;
                canvasHeight *= safeArea.height / screenHeight;

                canvasBottomLeftX = canvasRawSize.x * (safeArea.x / screenWidth);
                canvasBottomLeftY = canvasRawSize.y * (safeArea.y / screenHeight);
#endif
            }

            // Calculate safe area position of the popup
            // normalizedPosition allows us to glue the popup to a specific edge of the screen. It becomes useful when
            // the popup is at the right edge and we switch from portrait screen orientation to landscape screen orientation.
            // Without normalizedPosition, popup could jump to bottom or top edges instead of staying at the right edge
            var pos = canvasRawSize * 0.5f + (immediately
                ? new Vector2(normalizedPosition.x * canvasWidth, normalizedPosition.y * canvasHeight)
                : popupTransform.anchoredPosition - new Vector2(canvasBottomLeftX, canvasBottomLeftY));

            // Find distances to all four edges of the safe area
            var distToLeft = pos.x;
            var distToRight = canvasWidth - distToLeft;

            var distToBottom = pos.y;
            var distToTop = canvasHeight - distToBottom;

            var horDistance = Mathf.Min(distToLeft, distToRight);
            var vertDistance = Mathf.Min(distToBottom, distToTop);

            // Find the nearest edge's safe area coordinates
            if (horDistance < vertDistance)
            {
                if (distToLeft < distToRight)
                    pos = new Vector2(halfSize.x, pos.y);
                else
                    pos = new Vector2(canvasWidth - halfSize.x, pos.y);

                pos.y = Mathf.Clamp(pos.y, halfSize.y, canvasHeight - halfSize.y);
            }
            else
            {
                if (distToBottom < distToTop)
                    pos = new Vector2(pos.x, halfSize.y);
                else
                    pos = new Vector2(pos.x, canvasHeight - halfSize.y);

                pos.x = Mathf.Clamp(pos.x, halfSize.x, canvasWidth - halfSize.x);
            }

            pos -= canvasRawSize * 0.5f;

            normalizedPosition.Set(pos.x / canvasWidth, pos.y / canvasHeight);

            // Safe area's bottom left coordinates are added to pos only after normalizedPosition's value
            // is set because normalizedPosition is in range [-canvasWidth / 2, canvasWidth / 2]
            pos += new Vector2(canvasBottomLeftX, canvasBottomLeftY);

            // If another smooth movement animation is in progress, cancel it
            if (moveToPosCoroutine != null)
            {
                StopCoroutine(moveToPosCoroutine);
                moveToPosCoroutine = null;
            }

            if (immediately)
            {
                popupTransform.anchoredPosition = pos;
            }
            else
            {
                // Smoothly translate the popup to the specified position
                moveToPosCoroutine = MoveToPosAnimation(pos);
                StartCoroutine(moveToPosCoroutine);
            }
        }

#pragma warning disable 0649
        [SerializeField] private DebugLogManager debugManager;

        [SerializeField] private Text newInfoCountText;
        [SerializeField] private Text newWarningCountText;
        [SerializeField] private Text newErrorCountText;

        [SerializeField] private Color alertColorInfo;
        [SerializeField] private Color alertColorWarning;
        [SerializeField] private Color alertColorError;
#pragma warning restore 0649
    }
}