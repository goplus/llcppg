namespace bar {
	class base
	{
	public:
		~base();
	};

	class derived
    {
    private:
        base b;
    }
}
