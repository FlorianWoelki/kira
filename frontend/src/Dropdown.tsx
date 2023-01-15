import React from 'react';
import clsx from 'clsx';
import { ChevronDownIcon } from '@heroicons/react/20/solid';
import { Menu, Transition } from '@headlessui/react';

interface Props {
  title: string;
  items: string[];
  className?: string;
  onSelect?: (item?: string) => void;
}

export const Dropdown: React.FC<Props> = ({
  title,
  items,
  className,
  onSelect,
}): JSX.Element => {
  return (
    <Menu
      as="div"
      className={clsx('relative inline-block text-left', className)}
    >
      <div>
        <Menu.Button className="inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:ring-offset-gray-100">
          {title}
          <ChevronDownIcon className="-mr-1 ml-2 h-5 w-5" aria-hidden="true" />
        </Menu.Button>
      </div>

      <Transition
        as={React.Fragment}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
      >
        <Menu.Items className="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
          <div className="py-1">
            {items.map((item) => (
              <Menu.Item key={item}>
                {({ active }) => (
                  <a
                    onClick={() => onSelect?.(item)}
                    href="#"
                    className={clsx(
                      active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                      'block px-4 py-2 text-sm',
                    )}
                  >
                    {item}
                  </a>
                )}
              </Menu.Item>
            ))}
          </div>
        </Menu.Items>
      </Transition>
    </Menu>
  );
};
