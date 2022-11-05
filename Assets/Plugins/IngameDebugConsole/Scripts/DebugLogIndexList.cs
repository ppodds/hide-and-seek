using System;

namespace IngameDebugConsole
{
    public class DebugLogIndexList<T>
    {
        private T[] indices;

        public DebugLogIndexList()
        {
            indices = new T[64];
            Count = 0;
        }

        public int Count { get; private set; }

        public T this[int index]
        {
            get => indices[index];
            set => indices[index] = value;
        }

        public void Add(T value)
        {
            if (Count == indices.Length)
                Array.Resize(ref indices, Count * 2);

            indices[Count++] = value;
        }

        public void Clear()
        {
            Count = 0;
        }

        public int IndexOf(T value)
        {
            return Array.IndexOf(indices, value);
        }
    }
}